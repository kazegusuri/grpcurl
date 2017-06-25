package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/spf13/pflag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

var (
	insecure = pflag.BoolP("insecure", "k", false, "with insecure")
	headers  = pflag.StringArrayP("header", "H", nil, "")
)

func listServices(rcli *grpcreflect.Client) error {
	svcs, err := rcli.ListServices()
	if err != nil {
		return err
	}

	for i := range svcs {
		fmt.Printf("%s\n", svcs[i])
	}

	return nil
}

func listMethods(rcli *grpcreflect.Client, serviceName string) error {
	sdesc, err := rcli.ResolveService(serviceName)
	if err != nil {
		return err
	}

	for _, mdesc := range sdesc.GetMethods() {
		fmt.Printf("%s\n", mdesc.GetName())
	}

	return nil
}

func listMethodsDetails(rcli *grpcreflect.Client, serviceName string) error {
	sdesc, err := rcli.ResolveService(serviceName)
	if err != nil {
		return err
	}

	for _, mdesc := range sdesc.GetMethods() {
		fmt.Printf("%s\n", mdesc.GetFullyQualifiedName())
		inType := mdesc.GetInputType()
		outType := mdesc.GetOutputType()
		inRPCType := ""
		outRPCType := ""
		if mdesc.IsClientStreaming() {
			inRPCType = "streaming "
		}
		if mdesc.IsServerStreaming() {
			outRPCType = "streaming "
		}
		fmt.Printf("  IN : %s%s\n  OUT: %s%s\n", inRPCType, inType.GetFullyQualifiedName(), outRPCType, outType.GetFullyQualifiedName())
	}

	return nil
}

func call(ctx context.Context, rcli *grpcreflect.Client, stub grpcdynamic.Stub, serviceName, methodName string, reader io.Reader, verbose bool) error {
	sdesc, err := rcli.ResolveService(serviceName)
	if err != nil {
		return fmt.Errorf("service couldn't be resolve: %v", err)
	}

	mdesc := sdesc.FindMethodByName(methodName)
	if mdesc == nil {
		return fmt.Errorf("method couldn't be found")
	}

	msg, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to ReadAll %v", err)
	}

	reqMsg := dynamic.NewMessage(mdesc.GetInputType())
	if err = reqMsg.UnmarshalJSON(msg); err != nil {
		return fmt.Errorf("unmarshal %v", err)
	}

	reqJSON, err := reqMsg.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal %v", err)
	}

	var headerMD metadata.MD
	var trailerMD metadata.MD
	resp, err := stub.InvokeRpc(ctx, mdesc, reqMsg, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("unknown error: %v", err)
		}

		resp = st.Proto()
	}

	marshaler := &jsonpb.Marshaler{}
	respJSON, err := marshaler.MarshalToString(resp)
	if err != nil {
		return fmt.Errorf("marshal %v", err)
	}

	if verbose {
		fmt.Printf("Request Message\n")
		fmt.Printf("%s\n", string(reqJSON))
		fmt.Printf("Response Message\n")
	}
	fmt.Printf("%s\n", respJSON)
	if verbose {
		fmt.Printf("Response Headers\n")
		for k, vs := range headerMD {
			for i := range vs {
				fmt.Printf("%s: %s\n", k, vs[i])
			}
		}

		fmt.Printf("Response Trailer\n")
		for k, vs := range trailerMD {
			for i := range vs {
				fmt.Printf("%s: %s\n", k, vs[i])
			}
		}
	}

	return nil
}

func showHelp() {
	fmt.Fprintf(os.Stderr, "try grpcurl --help")
	os.Exit(1)
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

func main() {
	pflag.Parse()

	args := pflag.Args()
	if len(args) < 2 {
		showHelp()
	}

	operation := args[0]
	switch operation {
	case "call":
	case "list_services":
	case "list_methods":
	case "list_methods_details":
	default:
		showHelp()
	}

	addr := args[1]
	if addr == "" {
		showHelp()
	}

	var dialOpts []grpc.DialOption
	if *insecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	cli := rpb.NewServerReflectionClient(conn)
	rcli := grpcreflect.NewClient(ctx, cli)

	outgoingMD := buildOutgoingMetadata(*headers)
	ctx = metadata.NewOutgoingContext(ctx, outgoingMD)

	switch operation {
	case "call":
		if len(args) != 4 {
			showHelp()
		}
		serviceName := args[2]
		methodName := args[3]
		stub := grpcdynamic.NewStub(conn)
		if err = call(ctx, rcli, stub, serviceName, methodName, os.Stdin, true); err != nil {
			fmt.Fprintf(os.Stderr, "call error: %v", err)
			os.Exit(1)
		}
	case "list_services":
		listServices(rcli)
	case "list_methods":
		if len(args) != 3 {
			showHelp()
		}
		serviceName := args[2]
		listMethods(rcli, serviceName)
	case "list_methods_details":
		if len(args) != 3 {
			showHelp()
		}
		serviceName := args[2]
		listMethodsDetails(rcli, serviceName)
	}
}
