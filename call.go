package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type CallCommand struct {
	cmd         *cobra.Command
	opts        *GlobalOptions
	headers     []string
	addr        string
	rcli        *grpcreflect.Client
	stub        grpcdynamic.Stub
	marshaler   *jsonpb.Marshaler
	unmarshaler *jsonpb.Unmarshaler
}

func NewCallCommand(opts *GlobalOptions) *CallCommand {
	c := &CallCommand{
		cmd: &cobra.Command{
			Use:   "call ADDR FULL_METHOD_NAME",
			Short: "Call gRPC method with JSON",
			Example: `
* call
echo '{"message": "hello"}' | grpcurl call localhost:8888 test.Test.Echo
`,
			Args:         cobra.ExactArgs(2),
			SilenceUsage: true,
		},
		opts: opts,
	}
	c.cmd.RunE = c.Run
	c.cmd.Flags().StringArrayVarP(&c.headers, "header", "H", nil, "header")
	return c
}

func (c *CallCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *CallCommand) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	c.addr = args[0]
	conn, err := NewGRPCConnection(ctx, c.addr, c.opts.Insecure)
	if err != nil {
		return err
	}
	defer conn.Close()
	c.rcli = NewServerReflectionClient(ctx, conn)
	c.stub = grpcdynamic.NewStub(conn)
	c.marshaler = &jsonpb.Marshaler{
		OrigName:     true,
		EmitDefaults: true,
	}
	c.unmarshaler = &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	if err := c.call(ctx, args[1], c.opts.Input); err != nil {
		return err
	}
	return nil
}

func buildOutgoingMetadata(header []string) metadata.MD {
	var pairs []string
	for i := range header {
		parts := strings.SplitN(header[i], ":", 2)
		if len(parts) < 2 {
			// todo: logging?
			continue
		}

		k, v := strings.TrimLeft(parts[0], " "), strings.TrimLeft(parts[1], " ")
		pairs = append(pairs, k, v)
	}
	return metadata.Pairs(pairs...)
}

func (c CallCommand) resolveMessage(fullMethodName string) (*desc.MethodDescriptor, error) {
	// assume that fully-qualified method name cosists of
	// FULL_SERVER_NAME + "." + METHOD_NAME
	// so split the last dot to get service name
	n := strings.LastIndex(fullMethodName, ".")
	if n < 0 {
		return nil, fmt.Errorf("invalid method name: %v", fullMethodName)
	}
	serviceName := fullMethodName[0:n]
	methodName := fullMethodName[n+1:]

	sdesc, err := c.rcli.ResolveService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("service couldn't be resolve: %v", err)
	}

	mdesc := sdesc.FindMethodByName(methodName)
	if mdesc == nil {
		return nil, fmt.Errorf("method couldn't be found")
	}

	return mdesc, nil
}

func (c CallCommand) createMessage(mdesc *desc.MethodDescriptor, r io.Reader) (*dynamic.Message, error) {
	msg := dynamic.NewMessage(mdesc.GetInputType())
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to ReadAll %v", err)
	}
	if err = msg.UnmarshalJSONPB(c.unmarshaler, input); err != nil {
		return nil, fmt.Errorf("unmarshal %v", err)
	}
	return msg, nil
}

func (c CallCommand) call(ctx context.Context, fullMethodName string, reader io.Reader) error {
	mdesc, err := c.resolveMessage(fullMethodName)
	if err != nil {
		return err
	}

	msg, err := c.createMessage(mdesc, reader)
	if err != nil {
		return err
	}

	if c.opts.Verbose {
		reqJSON, err := msg.MarshalJSONPB(c.marshaler)
		if err != nil {
			return fmt.Errorf("marshal %v", err)
		}
		fmt.Printf("==> Request Message\n")
		fmt.Printf("%s\n", string(reqJSON))
	}

	var headerMD metadata.MD
	var trailerMD metadata.MD
	resp, err := c.stub.InvokeRpc(ctx, mdesc, msg, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("unknown error: %v", err)
		}

		resp = st.Proto()
	}

	respJSON, err := c.marshaler.MarshalToString(resp)
	if err != nil {
		return fmt.Errorf("marshal %v", err)
	}

	if c.opts.Verbose {
		fmt.Printf("<== Response Message\n")
	}
	fmt.Printf("%s\n", respJSON)
	if c.opts.Verbose {
		fmt.Printf("<== Response Headers\n")
		for k, vs := range headerMD {
			for i := range vs {
				fmt.Printf("%s: %s\n", k, vs[i])
			}
		}

		fmt.Printf("<== Response Trailer\n")
		for k, vs := range trailerMD {
			for i := range vs {
				fmt.Printf("%s: %s\n", k, vs[i])
			}
		}
	}

	return nil
}
