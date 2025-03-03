package codegen

import (
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	v3helper "github.com/unionj-cloud/go-doudou/cmd/internal/openapi/v3"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"os"
	"path/filepath"
	"testing"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testdata")
}

func TestGenDoc(t *testing.T) {
	dir := testDir + "doc1"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	svcfile := filepath.Join(dir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, ExprStringP)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				dir,
				ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenDoc(tt.args.dir, tt.args.ic, 1)
		})
	}
}

func TestGenDocUploadFile(t *testing.T) {
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	svcfile := filepath.Join(testDir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, ExprStringP)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				testDir,
				ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenDoc(tt.args.dir, tt.args.ic, 1)
		})
	}
}

func Test_schemasOf(t *testing.T) {
	type args struct {
		vofile string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				vofile: pathutils.Abs("testdata") + "/vo/vo.go",
			},
			want: 6,
		},
		{
			name: "",
			args: args{
				vofile: pathutils.Abs("testdata") + "/vo/vo1.go",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v3helper.SchemaNames = getSchemaNames(tt.args.vofile)
			if got := schemasOf(tt.args.vofile); len(got) != tt.want {
				t.Errorf("schemasOf() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestParseVo(t *testing.T) {
	Convey("Test ParseVo", t, func() {
		So(func() {
			ParseVo(testDir)
		}, ShouldNotPanic)
		So(len(v3helper.Schemas), ShouldNotBeZeroValue)
	})
}

func TestParseVoPanicWalk(t *testing.T) {
	Convey("Test ParseVo should panic from Walk", t, func() {
		Stat = os.Stat
		Walk = filepath.Walk
		Walk = func(root string, walkFn filepath.WalkFunc) error {
			return errors.New("mock Walk error")
		}
		So(func() {
			ParseVo(testDir)
		}, ShouldPanic)
	})
}
