// https://cloud.google.com/appengine/docs/standard/go111/googlecloudstorageclient/read-write-to-cloud-storage
// https://github.com/GoogleCloudPlatform/golang-samples/blob/8deb2909eadf32523007fd8fe9e8755a12c6d463/docs/appengine/storage/app.go
// Google Cloud Storage API.
// Access control lists (ACLs) control access to the buckets and to the objects contained in them.

// https://medium.com/wesionary-team/golang-image-upload-with-google-cloud-storage-and-gin-part-1-e5e668c1a5e2
package google

import (
	"bytes"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type GoogleCloud struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle

	writer io.Writer
	ctx    context.Context
	// cleanUp is a list of filenames that need cleaning up at the end of the demo.
	cleanUp []string
	// failed indicates that one or more of the demo steps failed.
	failed bool
}

func (gc *GoogleCloud) errorf(format string, args ...interface{}) {
	gc.failed = true
	fmt.Fprintln(gc.writer, fmt.Sprintf(format, args...))
	log.Errorf(gc.ctx, format, args...)
}

// handler is the main demo entry point that calls the GCS operations.
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//[START get_default_bucket]
	// Use `dev_appserver.py --default_gcs_bucket_name GCS_BUCKET_NAME`
	// when running locally.
	/*
		bucket, err := file.DefaultBucketName(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
		}
	*/
	//[END get_default_bucket]

	//ctx := context.Background()

	// create a client
	//storageClient, err := storage.NewClient(ctx)
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("tdg-analytics-prasan-serviceaccount.json"))

	bucketName := "5ps-consent-truegroup-true-test"

	if err != nil {
		log.Errorf(ctx, "failed to create storageClient: %v", err)
		return
	}
	defer storageClient.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Demo GCS Application running from Version: %v\n", appengine.VersionID(ctx))
	fmt.Fprintf(w, "Using bucket name: %v\n\n", bucketName)

	buf := &bytes.Buffer{}

	gc := &GoogleCloud{
		writer:     buf,
		ctx:        ctx,
		client:     storageClient,
		bucket:     storageClient.Bucket(bucketName),
		bucketName: bucketName,
	}

	gc.listBucket()

	/*
		fileName := "demo-testfile-go"
		gc.createFile(fileName)
		gc.readFile(fileName)
		gc.copyFile(fileName)
		gc.statFile(fileName)
		gc.createListFiles()
		gc.listBucket()
		gc.listBucketDirMode()
		gc.defaultACL()
		gc.putDefaultACLRule()
		gc.deleteDefaultACLRule()
		gc.bucketACL()
		gc.putBucketACLRule()
		gc.deleteBucketACLRule()
		gc.acl(fileName)
		gc.putACLRule(fileName)
		gc.deleteACLRule(fileName)
		gc.deleteFiles()
	*/
	if gc.failed {
		w.WriteHeader(http.StatusInternalServerError)
		buf.WriteTo(w)
		fmt.Fprintf(w, "\nDemo failed.\n")
	} else {
		w.WriteHeader(http.StatusOK)
		buf.WriteTo(w)
		fmt.Fprintf(w, "\nDemo succeeded.\n")
	}
}

//[START write]
// createFile creates a file in Google Cloud Storage.
func (gc *GoogleCloud) createFile(fileName string) {
	fmt.Fprintf(gc.writer, "Creating file /%v/%v\n", gc.bucketName, fileName)

	wc := gc.bucket.Object(fileName).NewWriter(gc.ctx)
	wc.ContentType = "text/plain"
	wc.Metadata = map[string]string{
		"x-goog-meta-foo": "foo",
		"x-goog-meta-bar": "bar",
	}
	gc.cleanUp = append(gc.cleanUp, fileName)

	if _, err := wc.Write([]byte("abcde\n")); err != nil {
		gc.errorf("createFile: unable to write data to bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}
	if _, err := wc.Write([]byte(strings.Repeat("f", 1024*4) + "\n")); err != nil {
		gc.errorf("createFile: unable to write data to bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		gc.errorf("createFile: unable to close bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}
}

//[END write]

//[START read]
// readFile reads the named file in Google Cloud Storage.
func (gc *GoogleCloud) readFile(fileName string) {
	io.WriteString(gc.writer, "\nAbbreviated file content (first line and last 1K):\n")

	rc, err := gc.bucket.Object(fileName).NewReader(gc.ctx)
	if err != nil {
		gc.errorf("readFile: unable to open file from bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}
	defer rc.Close()
	slurp, err := io.ReadAll(rc)
	if err != nil {
		gc.errorf("readFile: unable to read data from bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}

	fmt.Fprintf(gc.writer, "%s\n", bytes.SplitN(slurp, []byte("\n"), 2)[0])
	if len(slurp) > 1024 {
		fmt.Fprintf(gc.writer, "...%s\n", slurp[len(slurp)-1024:])
	} else {
		fmt.Fprintf(gc.writer, "%s\n", slurp)
	}
}

//[END read]

//[START copy]
// copyFile copies a file in Google Cloud Storage.
func (gc *GoogleCloud) copyFile(fileName string) {
	copyName := fileName + "-copy"
	fmt.Fprintf(gc.writer, "Copying file /%v/%v to /%v/%v:\n", gc.bucketName, fileName, gc.bucketName, copyName)

	obj, err := gc.bucket.Object(copyName).CopierFrom(gc.bucket.Object(fileName)).Run(gc.ctx)
	if err != nil {
		gc.errorf("copyFile: unable to copy /%v/%v to bucket %q, file %q: %v", gc.bucketName, fileName, gc.bucketName, copyName, err)
		return
	}
	gc.cleanUp = append(gc.cleanUp, copyName)

	gc.dumpStats(obj)
}

//[END copy]

func (gc *GoogleCloud) dumpStats(obj *storage.ObjectAttrs) {
	fmt.Fprintf(gc.writer, "(filename: /%v/%v, ", obj.Bucket, obj.Name)
	fmt.Fprintf(gc.writer, "ContentType: %q, ", obj.ContentType)
	fmt.Fprintf(gc.writer, "ACL: %#v, ", obj.ACL)
	fmt.Fprintf(gc.writer, "Owner: %v, ", obj.Owner)
	fmt.Fprintf(gc.writer, "ContentEncoding: %q, ", obj.ContentEncoding)
	fmt.Fprintf(gc.writer, "Size: %v, ", obj.Size)
	fmt.Fprintf(gc.writer, "MD5: %q, ", obj.MD5)
	fmt.Fprintf(gc.writer, "CRC32C: %q, ", obj.CRC32C)
	fmt.Fprintf(gc.writer, "Metadata: %#v, ", obj.Metadata)
	fmt.Fprintf(gc.writer, "MediaLink: %q, ", obj.MediaLink)
	fmt.Fprintf(gc.writer, "StorageClass: %q, ", obj.StorageClass)
	if !obj.Deleted.IsZero() {
		fmt.Fprintf(gc.writer, "Deleted: %v, ", obj.Deleted)
	}
	fmt.Fprintf(gc.writer, "Updated: %v)\n", obj.Updated)
}

//[START file_metadata]
// statFile reads the stats of the named file in Google Cloud Storage.
func (gc *GoogleCloud) statFile(fileName string) {
	io.WriteString(gc.writer, "\nFile stat:\n")

	obj, err := gc.bucket.Object(fileName).Attrs(gc.ctx)
	if err != nil {
		gc.errorf("statFile: unable to stat file from bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}

	gc.dumpStats(obj)
}

//[END file_metadata]

// createListFiles creates files that will be used by listBucket.
func (gc *GoogleCloud) createListFiles() {
	io.WriteString(gc.writer, "\nCreating more files for listbucket...\n")
	for _, n := range []string{"foo1", "foo2", "bar", "bar/1", "bar/2", "boo/"} {
		gc.createFile(n)
	}
}

//[START list_bucket]
// listBucket lists the contents of a bucket in Google Cloud Storage.
func (gc *GoogleCloud) listBucket() {
	io.WriteString(gc.writer, "\nListbucket result:\n")

	query := &storage.Query{Prefix: "/inbound/20220515/*"}
	iter := gc.bucket.Objects(gc.ctx, query)

	for {
		obj, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			gc.errorf("listBucket: unable to list bucket %q: %v", gc.bucketName, err)
			return
		}
		gc.dumpStats(obj)
	}
}

//[END list_bucket]

func (gc *GoogleCloud) listDir(name, indent string) {
	query := &storage.Query{Prefix: name, Delimiter: "/"}
	it := gc.bucket.Objects(gc.ctx, query)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			gc.errorf("listBucketDirMode: unable to list bucket %q: %v", gc.bucketName, err)
			return
		}
		if obj.Prefix == "" {
			fmt.Fprint(gc.writer, indent)
			gc.dumpStats(obj)
			continue
		}
		fmt.Fprintf(gc.writer, "%v(directory: /%v/%v)\n", indent, gc.bucketName, obj.Prefix)
		gc.listDir(obj.Prefix, indent+"  ")
	}
}

// listBucketDirMode lists the contents of a bucket in dir mode in Google Cloud Storage.
func (gc *GoogleCloud) listBucketDirMode() {
	io.WriteString(gc.writer, "\nListbucket directory mode result:\n")
	gc.listDir("b", "")
}

// dumpDefaultACL prints out the default object ACL for this bucket.
func (gc *GoogleCloud) dumpDefaultACL() {
	acl, err := gc.bucket.ACL().List(gc.ctx)
	if err != nil {
		gc.errorf("defaultACL: unable to list default object ACL for bucket %q: %v", gc.bucketName, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(gc.writer, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

// defaultACL displays the default object ACL for this bucket.
func (gc *GoogleCloud) defaultACL() {
	io.WriteString(gc.writer, "\nDefault object ACL:\n")
	gc.dumpDefaultACL()
}

// putDefaultACLRule adds the "allUsers" default object ACL rule for this bucket.
func (gc *GoogleCloud) putDefaultACLRule() {
	io.WriteString(gc.writer, "\nPut Default object ACL Rule:\n")
	err := gc.bucket.DefaultObjectACL().Set(gc.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		gc.errorf("putDefaultACLRule: unable to save default object ACL rule for bucket %q: %v", gc.bucketName, err)
		return
	}
	gc.dumpDefaultACL()
}

// deleteDefaultACLRule deleted the "allUsers" default object ACL rule for this bucket.
func (gc *GoogleCloud) deleteDefaultACLRule() {
	io.WriteString(gc.writer, "\nDelete Default object ACL Rule:\n")
	err := gc.bucket.DefaultObjectACL().Delete(gc.ctx, storage.AllUsers)
	if err != nil {
		gc.errorf("deleteDefaultACLRule: unable to delete default object ACL rule for bucket %q: %v", gc.bucketName, err)
		return
	}
	gc.dumpDefaultACL()
}

// dumpBucketACL prints out the bucket ACL.
func (gc *GoogleCloud) dumpBucketACL() {
	acl, err := gc.bucket.ACL().List(gc.ctx)
	if err != nil {
		gc.errorf("dumpBucketACL: unable to list bucket ACL for bucket %q: %v", gc.bucketName, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(gc.writer, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

// bucketACL displays the bucket ACL for this bucket.
func (gc *GoogleCloud) bucketACL() {
	io.WriteString(gc.writer, "\nBucket ACL:\n")
	gc.dumpBucketACL()
}

// putBucketACLRule adds the "allUsers" bucket ACL rule for this bucket.
func (gc *GoogleCloud) putBucketACLRule() {
	io.WriteString(gc.writer, "\nPut Bucket ACL Rule:\n")
	err := gc.bucket.ACL().Set(gc.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		gc.errorf("putBucketACLRule: unable to save bucket ACL rule for bucket %q: %v", gc.bucketName, err)
		return
	}
	gc.dumpBucketACL()
}

// deleteBucketACLRule deleted the "allUsers" bucket ACL rule for this bucket.
func (gc *GoogleCloud) deleteBucketACLRule() {
	io.WriteString(gc.writer, "\nDelete Bucket ACL Rule:\n")
	err := gc.bucket.ACL().Delete(gc.ctx, storage.AllUsers)
	if err != nil {
		gc.errorf("deleteBucketACLRule: unable to delete bucket ACL rule for bucket %q: %v", gc.bucketName, err)
		return
	}
	gc.dumpBucketACL()
}

// dumpACL prints out the ACL of the named file.
func (gc *GoogleCloud) dumpACL(fileName string) {
	acl, err := gc.bucket.Object(fileName).ACL().List(gc.ctx)
	if err != nil {
		gc.errorf("dumpACL: unable to list file ACL for bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}
	for _, v := range acl {
		fmt.Fprintf(gc.writer, "Scope: %q, Permission: %q\n", v.Entity, v.Role)
	}
}

// acl displays the ACL for the named file.
func (gc *GoogleCloud) acl(fileName string) {
	fmt.Fprintf(gc.writer, "\nACL for file %v:\n", fileName)
	gc.dumpACL(fileName)
}

// putACLRule adds the "allUsers" ACL rule for the named file.
func (gc *GoogleCloud) putACLRule(fileName string) {
	fmt.Fprintf(gc.writer, "\nPut ACL rule for file %v:\n", fileName)
	err := gc.bucket.Object(fileName).ACL().Set(gc.ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		gc.errorf("putACLRule: unable to save ACL rule for bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}
	gc.dumpACL(fileName)
}

// deleteACLRule deleted the "allUsers" ACL rule for the named file.
func (gc *GoogleCloud) deleteACLRule(fileName string) {
	fmt.Fprintf(gc.writer, "\nDelete ACL rule for file %v:\n", fileName)
	err := gc.bucket.Object(fileName).ACL().Delete(gc.ctx, storage.AllUsers)
	if err != nil {
		gc.errorf("deleteACLRule: unable to delete ACL rule for bucket %q, file %q: %v", gc.bucketName, fileName, err)
		return
	}
	gc.dumpACL(fileName)
}

// deleteFiles deletes all the temporary files from a bucket created by this demo.
func (gc *GoogleCloud) deleteFiles() {
	io.WriteString(gc.writer, "\nDeleting files...\n")
	for _, v := range gc.cleanUp {
		fmt.Fprintf(gc.writer, "Deleting file %v\n", v)
		if err := gc.bucket.Object(v).Delete(gc.ctx); err != nil {
			gc.errorf("deleteFiles: unable to delete bucket %q, file %q: %v", gc.bucketName, v, err)
			return
		}
	}
}

func StartHttp() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
