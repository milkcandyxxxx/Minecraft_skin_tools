package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

//package main

var Logo = `                  ███  ████  █████                                        █████           
                 ▒▒▒  ▒▒███ ▒▒███                                        ▒▒███            
 █████████████   ████  ▒███  ▒███ █████  ██████   ██████   ████████    ███████  █████ ████
▒▒███▒▒███▒▒███ ▒▒███  ▒███  ▒███▒▒███  ███▒▒███ ▒▒▒▒▒███ ▒▒███▒▒███  ███▒▒███ ▒▒███ ▒███ 
 ▒███ ▒███ ▒███  ▒███  ▒███  ▒██████▒  ▒███ ▒▒▒   ███████  ▒███ ▒███ ▒███ ▒███  ▒███ ▒███ 
 ▒███ ▒███ ▒███  ▒███  ▒███  ▒███▒▒███ ▒███  ███ ███▒▒███  ▒███ ▒███ ▒███ ▒███  ▒███ ▒███ 
 █████▒███ █████ █████ █████ ████ █████▒▒██████ ▒▒████████ ████ █████▒▒████████ ▒▒███████ 
▒▒▒▒▒ ▒▒▒ ▒▒▒▒▒ ▒▒▒▒▒ ▒▒▒▒▒ ▒▒▒▒ ▒▒▒▒▒  ▒▒▒▒▒▒   ▒▒▒▒▒▒▒▒ ▒▒▒▒ ▒▒▒▒▒  ▒▒▒▒▒▒▒▒   ▒▒▒▒▒███ 
                                                                                 ███ ▒███ 
                                                                                ▒▒██████  
                                                                                 ▒▒▒▒▒▒   `

func main() {
	begin()
}
func begin() {
	fmt.Println(Logo)
	decode_file, err := openfile("xmilkcandy.png")
	if err != nil {
		fmt.Println(err)
	}
	ww(decode_file, 0, 0, 1, 0, 8, "xmilkcandy.png")

}

func openfile(file_name string) (image.Image, error) {
	photo_file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("打开图片失败")
		return nil, err
	}
	defer func(photo_file *os.File) {
		err := photo_file.Close()
		if err != nil {
		}
	}(photo_file)
	decode_file, _, err := image.Decode(photo_file)
	if err != nil {
		fmt.Println("解码图片失败")
		return nil, err
	}
	//defer photo_file.Close()
	//解码图片
	return decode_file, nil
}

func ww(decode_file image.Image, x_old int, y_old int, x_new int, y_new int, xy int, file_name string) error {
	file_max_x := decode_file.Bounds().Max.X - decode_file.Bounds().Min.X
	file_max_y := decode_file.Bounds().Max.Y - decode_file.Bounds().Min.Y
	if file_max_x != 64 && file_max_y != 64 {
		fmt.Println("这不是一个我的世界皮肤文件哦")
		return errors.New("这不是一个我的世界皮肤文件哦")
	}
	//复制一份图片的rgb便有修改
	newfile := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			newfile.Set(x, y, decode_file.At(x, y))
		}
	}
	//创建数组存储rgb数据[x][y][r,g,b,a]

	var file_temp [8][8][4]uint8
	for x := 0; x < xy; x++ {
		for y := 0; y < xy; y++ {
			r, g, b, a := newfile.At(xy*x_old+x, xy*y_old+y).RGBA()
			file_temp[x][y][0] = uint8(r >> 8)
			file_temp[x][y][1] = uint8(g >> 8)
			file_temp[x][y][2] = uint8(b >> 8)
			file_temp[x][y][3] = uint8(a >> 8)
		}
	}
	for x := 0; x < xy; x++ {
		for y := 0; y < xy; y++ {
			newfile.Set(xy*x_old+x, xy*y_old+y, decode_file.At(xy*x_new+x, xy*x_new+y))
		}
	}
	for x := 0; x < xy; x++ {
		for y := 0; y < xy; y++ {
			r, g, b, a := file_temp[x][y][0], file_temp[x][y][1], file_temp[x][y][2], file_temp[x][y][3]
			newfile.SetRGBA(xy*x_new+x, xy*x_new+y, color.RGBA{r, g, b, a})
		}
	}
	outfile, err := os.Create("new" + file_name)
	if err != nil {
		fmt.Println("创建图片失败")
		return err
	}
	defer outfile.Close()
	png.Encode(outfile, newfile)
	return nil

}
