package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/linyows/ogenexample/oas/oasgen"
)

func run(ctx context.Context) error {
	var arg struct {
		BaseURL string
		ID      int64
	}
	flag.StringVar(&arg.BaseURL, "url", "http://localhost:8080", "target server url")
	flag.Int64Var(&arg.ID, "id", 1, "pet id to request")
	flag.Parse()

	client, err := oasgen.NewClient(arg.BaseURL)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	res, err := client.GetPetById(ctx, oasgen.GetPetByIdParams{
		PetId: arg.ID,
	})
	if err != nil {
		return fmt.Errorf("get pet %d: %w", arg.ID, err)
	}

	switch p := res.(type) {
	case *oasgen.Pet:
		data, err := p.MarshalJSON()
		if err != nil {
			return err
		}
		var out bytes.Buffer
		if err := json.Indent(&out, data, "", "  "); err != nil {
			return err
		}
		color.New(color.FgGreen).Println(out.String())
	case *oasgen.GetPetByIdNotFound:
		return errors.New("not found")
	}

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		color.New(color.FgRed).Printf("%+v\n", err)
		os.Exit(2)
	}
}
