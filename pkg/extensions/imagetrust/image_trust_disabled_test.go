//go:build !imagetrust

package imagetrust_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"zotregistry.io/zot/pkg/extensions/imagetrust"
	"zotregistry.io/zot/pkg/test"
)

func TestImageTrust(t *testing.T) {
	Convey("binary doesn't include imagetrust", t, func() {
		rootDir := t.TempDir()

		cosignDir := path.Join(rootDir, "_cosign")
		_, err := os.Stat(cosignDir)
		So(os.IsNotExist(err), ShouldBeTrue)

		notationDir := path.Join(rootDir, "_notation")
		_, err = os.Stat(notationDir)
		So(os.IsNotExist(err), ShouldBeTrue)

		repo := "repo"

		image, err := test.GetRandomImage() //nolint:staticcheck
		So(err, ShouldBeNil)

		manifestContent, err := json.Marshal(image.Manifest)
		So(err, ShouldBeNil)

		manifestDigest := image.Digest()

		localImgTrustStore, err := imagetrust.NewLocalImageTrustStore(rootDir)
		So(err, ShouldBeNil)

		author, expTime, ok, err := localImgTrustStore.VerifySignature("cosign",
			[]byte(""), "", manifestDigest, manifestContent, repo,
		)
		So(author, ShouldBeEmpty)
		So(expTime, ShouldBeZeroValue)
		So(ok, ShouldBeFalse)
		So(err, ShouldBeNil)

		_, err = os.Stat(cosignDir)
		So(os.IsNotExist(err), ShouldBeTrue)

		_, err = os.Stat(notationDir)
		So(os.IsNotExist(err), ShouldBeTrue)

		cloudImgTrustStore, err := imagetrust.NewAWSImageTrustStore("region",
			"endpoint",
		)
		So(err, ShouldBeNil)

		author, expTime, ok, err = cloudImgTrustStore.VerifySignature("cosign",
			[]byte(""), "", manifestDigest, manifestContent, repo,
		)
		So(author, ShouldBeEmpty)
		So(expTime, ShouldBeZeroValue)
		So(ok, ShouldBeFalse)
		So(err, ShouldBeNil)
	})
}