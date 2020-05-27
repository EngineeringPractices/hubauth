package main

import (
	"context"
	"log"
	"net/http"
	"os"

	google_datastore "cloud.google.com/go/datastore"
	kms "cloud.google.com/go/kms/apiv1"
	"github.com/flynn/hubauth/pkg/datastore"
	"github.com/flynn/hubauth/pkg/httpapi"
	"github.com/flynn/hubauth/pkg/idp"
	"github.com/flynn/hubauth/pkg/kmssign"
	"github.com/flynn/hubauth/pkg/rp/google"
)

func main() {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8000"
	}

	ctx := context.Background()
	dsClient, err := google_datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatalf("error initializing datastore client: %s", err)
	}

	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		log.Fatalf("error initializing kms client: %s", err)
	}
	cookieKey, err := kmssign.NewKey(ctx, kmsClient, os.Getenv("COOKIE_KEY"))
	if err != nil {
		log.Fatalf("error initializing cookie key: %s", err)
	}
	rpKey, err := kmssign.NewKey(ctx, kmsClient, os.Getenv("RP_KEY"))
	if err != nil {
		log.Fatalf("error initializing rp key: %s", err)
	}
	refreshKey, err := kmssign.NewKey(ctx, kmsClient, os.Getenv("REFRESH_KEY"))
	if err != nil {
		log.Fatalf("error initializing refresh key: %s", err)
	}
	accessKey, err := kmssign.NewKey(ctx, kmsClient, os.Getenv("ACCESS_KEY"))
	if err != nil {
		log.Fatalf("error initializing access key: %s", err)
	}

	log.Fatal(http.ListenAndServe(":"+httpPort, httpapi.New(
		idp.New(datastore.New(dsClient),
			google.New(
				os.Getenv("RP_GOOGLE_CLIENT_ID"),
				os.Getenv("RP_GOOGLE_CLIENT_SECRET"),
				os.Getenv("BASE_URL")+"/rp/google",
				rpKey,
			),
			refreshKey,
			accessKey,
		),
		cookieKey,
	)))
}