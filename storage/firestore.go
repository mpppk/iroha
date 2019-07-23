package storage

import (
	"context"
	"fmt"
	"reflect"
	"time"

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
	Progress int
}

func NewFireStore(ctx context.Context, filePath string, rootCollectionName string, projectId string) (_ *FireStore, err error) {
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
	bucketClient, err := newCloudStorageClient(ctx, app, projectId, rootCollectionName, bucketAttrs)
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

func (f *FireStore) getProgress(ctx context.Context, indices []int) (int, error) {
	doc, err := f.client.Collection(f.rootCollectionName).Doc(toStorageStrKey(indices)).Get(ctx)
	if doc != nil && !doc.Exists() {
		return progressNotStarted, nil
	}
	if err != nil {
		return -1, errors.Wrap(err, "failed to get progress from firestore")
	}
	m := doc.Data()
	progressInterface, ok := m["Progress"]
	if !ok {
		return -1, fmt.Errorf("failed to get progress from firestore. Progress key not found. indicdes: %v", indices)
	}
	progress, ok := progressInterface.(int64)
	if !ok {
		return -1, fmt.Errorf("failed to get progress from firestore. Invalid type value of Progress key. indicdes: %v v: %v(%s)",
			indices, progressInterface, reflect.TypeOf(progressInterface))
	}
	return int(progress), nil
}

func (f *FireStore) Start(ctx context.Context, indices []int) error {
	return f.updateProgress(ctx, indices, progressProcessing)
}

func (f *FireStore) updateProgress(ctx context.Context, indices []int, progress int) error {
	doc := &cacheDoc{Progress: progress}
	_, err := f.client.Collection(f.rootCollectionName).Doc(toStorageStrKey(indices)).Set(ctx, doc)
	if err != nil {
		return errors.Wrapf(err, "failed to update progress on firestore. indices: %s", indices)
	}
	return nil
}

func (f *FireStore) Set(ctx context.Context, indices []int, wordsList [][]*ktkn.Word) error {
	wl := wordsList
	if wl == nil {
		wl = make([][]*ktkn.Word, 0)
	}

	if err := f.cstorageClient.SaveWordsList(ctx, indices, wordsList); err != nil {
		return errors.Wrap(err, "failed to set wordsList to cloud storage")
	}
	if err := f.updateProgress(ctx, indices, progressDone); err != nil {
		return errors.Wrap(err, "failed to update progress to done")
	}

	return nil
}

func (f *FireStore) Get(ctx context.Context, indices []int) ([][]*ktkn.Word, bool, error) {
L:
	for {
		// FIXME: watch document
		progress, err := f.getProgress(ctx, indices)
		if err != nil {
			return nil, false, err
		}

		switch progress {
		case progressNotStarted:
			return nil, false, nil
		case progressDone:
			break L
		case progressProcessing:
			time.Sleep(500 * time.Millisecond)
		}
	}

	return f.cstorageClient.Get(ctx, indices)
}
