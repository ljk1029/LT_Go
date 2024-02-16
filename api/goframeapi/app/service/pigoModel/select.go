package pigoModel

import (
	"fmt"
	"github.com/disintegration/imageorient"
	"github.com/nfnt/resize"
	"github.com/gogf/gf/util/guid"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
)

func Crop(imgUrl string, width, height int, minWidth, minHeight uint) (code, x0, y0 int, fileName string) {
	file, err := http.Get(imgUrl)
	print(fmt.Sprintf("%s\n",imgUrl))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(imgUrl)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fileConfig, _, err := imageorient.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 解码图像
	img, err := jpeg.Decode(file.Body)
	if err != nil {
		log.Fatal(err)
	}
	//file.Body.Close()

	//使用Lanczos重采样将宽度调整为1000
	//并保持高宽比
	if fileConfig.Bounds().Dx() > fileConfig.Bounds().Dy() {
		minWidth = 0
	} else {
		minHeight = 0
	}
	m := resize.Resize(minWidth, minHeight, img, resize.Lanczos2)

	//file, err := http.Get(imgUrl)
	//if err != nil {
	//	fmt.Println(err)
	//}
	defer file.Body.Close()
	//img, err := jpeg.Decode(file.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	stretchFileName := guid.S();
	stretchFileSrc := fmt.Sprintf("/home/img/%s.jpg",stretchFileName)
	fmt.Println(stretchFileName)
	out, err := os.Create(stretchFileSrc)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	//将新图片写入文件
	jpeg.Encode(out, m, nil)

	//读取重新定义好尺寸的本地图片
	file1, err := os.Open(stretchFileSrc)
	if err != nil {
		fmt.Println(err)
	}
	defer file1.Close()
	img1, err := jpeg.Decode(file1)
	if err != nil {
		fmt.Println(err)
	}
	code, x, y := DetectFaceNew(img1, SaveNew)
	fmt.Println("888888888888888888888888888888888")
	fmt.Println(code)
	fmt.Println("888888888888888888888888888888888")
	fmt.Println(x)
	fmt.Println("888888888888888888888888888888888")
	fmt.Println(y)
	fmt.Println("888888888888888888888888888888888")
	fmt.Println(stretchFileSrc)
	fmt.Println("888888888888888888888888888888888")
	//根据尺寸计算左上角x,y
	return code, x - width/2, y - height/2, stretchFileSrc
}
