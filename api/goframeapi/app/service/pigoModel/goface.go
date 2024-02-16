package pigoModel

import (
	"fmt"
	gohash "github.com/corona10/goimagehash"
	pigo "github.com/esimov/pigo/core"
	"github.com/gogf/gf/frame/g"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"
)

// Threshold of getArr.
const Threshold = 3.14159

// Paths of facefinder.(split by ;)
//const Paths = "/home/facefinder"
// classifier is threadsafe.
var classifier *pigo.Pigo

//保存探测器的结果。
func save(src *image.NRGBA, dets []pigo.Detection) []image.Image {
	for i, v := range dets {
		x, y, w := v.Col, v.Row, v.Scale/2
		print(x - w)
		print("\n")
		print(y - w)
		print("\n")
		print(x + w)
		print("\n")
		print(y + w)
		print("\n")

		print("2222222\n")
		img := src.SubImage(image.Rect(x-w, y-w, x+w, y+w))
		file, err := os.Create(fmt.Sprintf("%d-%.2f", i, v.Q) + ".jpg")
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer file.Close()
		if err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

//获取人脸中心点
func SaveNew(src *image.NRGBA, dets []pigo.Detection) (code, q, e int) {
	//未检测到人像
	if len(dets) == 0 {
		return 500, 0, 0
	}

	//特殊图片只解析到一组数据
	if len(dets) == 1 {
		return 0, dets[0].Col, dets[0].Row
	}

	//判断哪一个元素是人像中中心点数据
	a := dets[0]
	b := dets[1]
	aw := a.Scale / 2
	bw := b.Scale / 2

	if bw > aw {
		return 0, b.Col, b.Row
	} else {
		return 0, a.Col, a.Row
	}

	//for i, v := range dets {
	//
	//	x, y := v.Col, v.Row
	//	if i==0 {
	//
	//		return 0,	x, y
	//	}
	//	//x, y, w := v.Col, v.Row, v.Scale/2
	//	//print(i)
	//	//print("\n")
	//	//print(x - w)
	//	//print("\n")
	//	//print(y - w)
	//	//print("\n")
	//	//print(x + w)
	//	//print("\n")
	//	//print(y + w)
	//	//print("\n")
	//	//print(fmt.Sprintf("x坐标:%d\n",x))
	//	//print(fmt.Sprintf("y坐标:%d\n",y))
	//	//print(fmt.Sprintf("宽度:%d\n",w))
	//	//print("2222222\n")
	//	//
	//	//img := src.SubImage(image.Rect(x-w, y-w, x+w, y+w))
	//	//file, err := os.Create(fmt.Sprintf("%d-%.2f", i, v.Q) + ".jpg")
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//	continue
	//	//}
	//	//defer file.Close()
	//	//if err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
	//	//	fmt.Println(err)
	//	//}
	//}
	return 0, 0, 0

}

// a plugin must contains main function.

// callback 回调处理逻辑.
type callback func(src *image.NRGBA, dets []pigo.Detection) []image.Image

type callbackNew func(src *image.NRGBA, dets []pigo.Detection) (code, x, y int)

//得到检测面的结果.
func getArr(src *image.NRGBA, dets []pigo.Detection) []image.Image {
	var r []image.Image
	for _, v := range dets {
		if v.Q > Threshold {
			x, y, w := v.Col, v.Row, v.Scale/2
			img := src.SubImage(image.Rect(x-w, y-w, x+w, y+w))
			r = append(r, img)
		}
	}
	print("33333\n")
	return r
}
//得到检测面的结果.
func GetArr(src *image.NRGBA, dets []pigo.Detection) []image.Image {
	var r []image.Image
	for _, v := range dets {
		if v.Q > Threshold {
			x, y, w := v.Col, v.Row, v.Scale/2
			img := src.SubImage(image.Rect(x-w, y-w, x+w, y+w))
			r = append(r, img)
		}
	}
	print("33333\n")
	return r
}
// readFinder from files.
func readFinder(files ...string) ([]byte, error) {
	var cascadeFile []byte
	var err error
	for _, v := range files {
		cascadeFile, err = ioutil.ReadFile(v)
		if err == nil {
			fmt.Printf("Read the cascade file succeed: %v\n", v)
			break
		}
		fmt.Printf("Error reading the cascade file. %v\n", err)
	}
	return cascadeFile, err
}

// init the classifier.
func init() {
	if classifier != nil {
		return
	}
	Paths := g.Cfg("config").GetString("IsHeadPg.Route")
	paths := strings.Split(Paths, ";")
	cascadeFile, err := readFinder(paths...)
	if err != nil {
		return
	}
	if len(cascadeFile) == 0 {
		fmt.Printf("Error reading the cascade file: Empty file.\n")
		return
	}

	pigo := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err = pigo.Unpack(cascadeFile)
	if err != nil {
		fmt.Printf("Error reading the cascade file. %s\n", err)
	}
}

// DetectFace in a picture.
func DetectFace(img image.Image, cb callback) []image.Image {
	if classifier == nil || img == nil {
		fmt.Printf("The classifier or image is nil\n")
		return nil
	}

	src := pigo.ImgToNRGBA(img)
	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     32,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	angle := 0.0 // cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)
	if cb != nil {
		return cb(src, dets)
	}
	return nil
}

func DetectFaceNew(img image.Image, cb callbackNew) (code, x, y int) {
	if classifier == nil || img == nil {
		//图像异常
		return 200, 0, 0
	}
	src := pigo.ImgToNRGBA(img)
	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     32,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	angle := 0.0 // cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)
	if cb != nil {
		return cb(src, dets)
	}
	return 200, 0, 0
}

// imageCompare 图片比对算法.
func imageCompare(src *gohash.ImageHash, cmp image.Image) float64 {
	if src != nil {
		hash, _ := gohash.AverageHash(cmp)
		if n, err := src.Distance(hash); err == nil {
			return 1 - float64(n)/64.0
		}
	}
	return 0
}

// AlarmProcess 告警处理单元.
// go build -buildmode=plugin goface.go
func AlarmProcess(dis map[string]interface{}, features []interface{}, arr []image.Image, ids []string, level int) bool {
	var levelThresholdMap = map[int]float64{0: 0.8, 1: 0.6, 2: 0.8, 3: 0.9}
	threshold := levelThresholdMap[level]
	if _, ok := dis["hash"]; !ok { // 特征计算
		if img := dis["image"]; img != nil {
			if v, ok := img.(image.Image); ok {
				if varr := DetectFace(v, getArr); len(varr) > 0 {
					if hash, err := gohash.AverageHash(varr[0]); err == nil {
						dis["hash"] = hash
					}
				} else {
					fmt.Printf("未从布控图像检测到人脸.\n")
				}
			}
		}
		fmt.Printf("计算布控图像的特征值=%v.\n", dis["hash"])
	}
	if hash, ok := dis["hash"]; ok { // 图片比对
		for i, a := range arr {
			v, _ := hash.(*gohash.ImageHash)
			for _, img := range DetectFace(a, getArr) {
				f := imageCompare(v, img)
				if f > threshold {
					fmt.Printf("[%s]相似度阈值%f, 触发告警.\n", ids[i], f)
					return true
				}
				fmt.Printf("[%s]相似度阈值%f, 未触发告警.\n", ids[i], f)
			}
		}
	}
	return false
}
