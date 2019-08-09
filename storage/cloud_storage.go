package storage

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"

	"firebase.google.com/go/storage"

	cstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/mpppk/iroha/ktkn"
	"github.com/pkg/errors"
)

type CloudStorageClient struct {
	client     *storage.Client
	bucketName string
	projectId  string
	attrs      *cstorage.BucketAttrs
}

func newCloudStorageClient(ctx context.Context, app *firebase.App, projectId string, bucketName string, attrs *cstorage.BucketAttrs) (*CloudStorageClient, error) {
	baseErrMsg := "failed to create new bucket client"
	client, err := app.Storage(ctx)
	if err != nil {
		return nil, errors.Wrap(err, baseErrMsg+"failed to create firebase storage client")
	}

	return &CloudStorageClient{
		client:     client,
		bucketName: bucketName,
		projectId:  projectId,
		attrs:      attrs,
	}, nil
}

func (b *CloudStorageClient) SaveWordsList(ctx context.Context, indices []int, wordsList [][]*ktkn.Word) (err error) {
	bucket, err := b.getBucket()
	if err != nil {
		return errors.Wrap(err, "failed to save bucket")
	}

	wc := bucket.Object(toStorageStrKey(indices)).NewWriter(ctx)
	defer func() {
		if err := wc.Close(); err != nil {
			err = errors.Wrap(err, "failed to close cloud storage bucket object")
		}
	}()
	wc.Metadata = map[string]string{
		"Content-Encoding": "gzip",
	}
	gw := gzip.NewWriter(wc)
	defer func() {
		if err := gw.Close(); err != nil {
			err = errors.Wrap(err, "failed to close gzip writer")
		}
	}()

	if err := json.NewEncoder(gw).Encode(wordsList); err != nil {
		return errors.Wrap(err, "failed to encode wordsList to json")
	}
	if err := gw.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush gzip writer")
	}
	return err
}

func (b *CloudStorageClient) Get(ctx context.Context, indices []int) ([][]*ktkn.Word, bool, error) {
	reader, ok, err := b.getObjectReader(ctx, indices)
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to get cache from cloud storage. indices: %v", indices)
	}
	if !ok {
		return nil, false, nil
	}

	gr, err := gzip.NewReader(reader)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to create gzip reader")
	}

	var wordsList [][]*ktkn.Word
	if err := json.NewDecoder(gr).Decode(&wordsList); err != nil {
		return nil, false, err
	}
	return wordsList, true, nil
}

func (b *CloudStorageClient) getObject(ctx context.Context, indices []int) (*cstorage.ObjectHandle, bool, error) {
	bucket, err := b.getBucket()
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get bucket of cloud storage")
	}
	object := bucket.Object(toStorageStrKey(indices))
	if _, err := object.Attrs(ctx); err != nil {
		return nil, false, nil
	}
	return object, true, nil
}

func (b *CloudStorageClient) getObjectWriter(ctx context.Context, indices []int) (io.Writer, bool, error) {
	object, ok, err := b.getObject(ctx, indices)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get object of cloud storage")
	}
	if !ok {
		return nil, false, nil
	}
	return object.NewWriter(ctx), ok, nil
}

func (b *CloudStorageClient) getObjectReader(ctx context.Context, indices []int) (io.Reader, bool, error) {
	object, ok, err := b.getObject(ctx, indices)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get object of cloud storage")
	}
	if !ok {
		return nil, false, nil
	}
	reader, err := object.NewReader(ctx)
	return reader, ok, err
}

func (b *CloudStorageClient) getBucket() (*cstorage.BucketHandle, error) {
	bucket, err := b.client.Bucket(b.bucketName)
	return bucket, errors.Wrap(err, "failed to handle bucket. bucket name: "+b.bucketName)
}

func (b *CloudStorageClient) createBucketIfDoesNotExist(ctx context.Context, bucketName string) (bool, error) {
	if b.isBucketExists(ctx, bucketName) {
		return false, nil
	}

	bucket, err := b.client.Bucket(bucketName)
	if err != nil {
		return false, errors.Wrap(err, "failed to handle bucket")
	}
	return true, bucket.Create(ctx, b.projectId, b.attrs)
}

func (b *CloudStorageClient) isBucketExists(ctx context.Context, bucketName string) bool {
	bucket, err := b.client.Bucket(bucketName)
	if err != nil {
		return false
	}

	_, err = bucket.Attrs(ctx)
	return err == nil
}
