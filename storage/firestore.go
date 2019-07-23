package storage

import (
	"context"

	cstorage "cloud.google.com/go/storage"

	"github.com/pkg/errors"

	"github.com/mpppk/iroha/ktkn"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FireStore struct {
	client             *firestore.Client
	cstorageClient     *CloudStorageClient
	rootCollectionName string
}

type cacheDoc struct {
	Progress string
}

func NewFireStore(ctx context.Context, filePath string, rootCollectionName string) (_ *FireStore, err error) {
	var baseErrMsg = "failed to create firestore storage: "
	sa := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, errors.Wrap(err, baseErrMsg+"failed to create firebase new app")
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, errors.Wrap(err, baseErrMsg+"failed to create firestore client")
	}

	// FIXME
	bucketAttrs := &cstorage.BucketAttrs{
		StorageClass: "STANDARD",
		Location:     "asia-northeast1",
	}
	// FIXME projectId
	bucketClient, err := newCloudStorageClient(ctx, app, "iroha-247312", rootCollectionName, bucketAttrs)
	if err != nil {
		return nil, errors.Wrap(err, baseErrMsg+"failed to create bucket client")
	}

	if _, err := bucketClient.createBucketIfDoesNotExist(ctx, rootCollectionName); err != nil {
		return nil, errors.Wrap(err, "failed to create new bucket to cloud storage. bucket name: "+rootCollectionName)
	}

	return &FireStore{
		client:             client,
		cstorageClient:     bucketClient,
		rootCollectionName: rootCollectionName,
	}, errors.Wrap(err, baseErrMsg+"failed to close firestore client")
}

func (f *FireStore) Set(ctx context.Context, indices []int, wordsList [][]*ktkn.Word) error {
	wl := wordsList
	if wl == nil {
		wl = make([][]*ktkn.Word, 0)
	}

	if err := f.cstorageClient.SaveWordsList(ctx, indices, wordsList); err != nil {
		return errors.Wrap(err, "failed to set wordsList to cloud storage")
	}

	doc := &cacheDoc{Progress: "done"}
	_, err := f.client.Collection(f.rootCollectionName).Doc(toStorageStrKey(indices)).Set(ctx, doc)
	if err != nil {
		return errors.Wrap(err, "failed to set cache to firestore")
	}
	return nil
}

func (f *FireStore) Get(ctx context.Context, indices []int) ([][]*ktkn.Word, bool, error) {
	return f.cstorageClient.Get(ctx, indices)
}
