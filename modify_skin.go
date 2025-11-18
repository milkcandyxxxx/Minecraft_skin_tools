package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// 生成图片
func generates_picture(file_name string, newfile *image.RGBA) error {
	outfile, err := os.Create("new" + file_name)
	if err != nil {
		fmt.Println("创建图片失败")
		return err
	}
	defer func(outfile *os.File) {
		_ = outfile.Close()
	}(outfile)
	_ = png.Encode(outfile, newfile)
	return nil
}

// 打开图片，解码图片，创建新图片
func openfile(file_name string) (image.Image, error, *image.RGBA) {
	photo_file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("打开图片失败")
		return nil, err, nil
	}
	defer func(photo_file *os.File) {
		err := photo_file.Close()
		if err != nil {
		}
	}(photo_file)
	decode_file, _, err := image.Decode(photo_file)
	if err != nil {
		fmt.Println("解码图片失败")
		return nil, err, nil
	}
	file_max_x := decode_file.Bounds().Max.X - decode_file.Bounds().Min.X
	file_max_y := decode_file.Bounds().Max.Y - decode_file.Bounds().Min.Y
	if file_max_x != 64 && file_max_y != 64 {
		fmt.Println("这不是一个我的世界皮肤文件哦")
		return nil, errors.New("这不是一个我的世界皮肤文件哦"), nil
	}
	newfile := image.NewRGBA(image.Rect(0, 0, 64, 64))
	// defer photo_file.Close()
	// 解码图片
	return decode_file, nil, newfile
}

// 根据出入值进行修改像素
func exchange(decode_file image.Image, x_old int, y_old int, x_new int, y_new int, xy int, file_name string, newfile *image.RGBA) (error, *image.RGBA) {

	// 复制一份图片的rgb便有修改

	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			newfile.Set(x, y, decode_file.At(x, y))
		}
	}
	// 创建数组存储rgb数据[x][y][r,g,b,a]

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
			newfile.Set(xy*x_old+x, xy*y_old+y, decode_file.At(xy*x_new+x, xy*y_new+y))
		}
	}
	for x := 0; x < xy; x++ {
		for y := 0; y < xy; y++ {
			r, g, b, a := file_temp[x][y][0], file_temp[x][y][1], file_temp[x][y][2], file_temp[x][y][3]
			newfile.SetRGBA(xy*x_new+x, xy*y_new+y, color.RGBA{r, g, b, a})
		}
	}

	return nil, newfile

}

// #########################################################################################预设############################################################################
// ########################################################################################################################################################################
// ########################################################################################################################################################################

func modify_skin() error {
	fmt.Printf(Logo2+Logo3+"%-30s"+Logo3+Logo2+"\n", "1.预设")
	fmt.Printf(Logo2+Logo3+"%-30s"+Logo3+Logo2+"\n", "2.自定义更换")
	fmt.Printf(Logo2+Logo3+"%-30s"+Logo3+Logo2+"\n", "3.退出")
	fmt.Print("请选择序号:")
	var choice int
	var choice_2 int
	fmt.Scan(&choice)
	for {
		switch choice {
		case 1:
			fmt.Print("输入文件的完整名字，且需要在当前目录（例xmilkcandy.png）")
			var file_name string
			fmt.Scanln(&file_name)
			fmt.Printf(Logo2+Logo3+"%-30s"+Logo3+Logo2+"\n", "1.头翻转")
			fmt.Printf(Logo2+Logo3+"%-30s"+Logo3+Logo2+"\n", "1.身体翻转")
			fmt.Printf(Logo2+Logo3+"%-30s"+Logo3+Logo2+"\n", "1.手和腿替换")
			fmt.Print("输入序号:")
			fmt.Scanln(&choice_2)
			switch choice_2 {
			case 1:
				err := change_heads(file_name)
				if err == nil {
					fmt.Println("按回车键退出...")
					fmt.Scanln()
					os.Exit(0)
				} else {
					continue
				}
			}

		case 2:
		case 3:
			fmt.Println("退出")
			os.Exit(0)
		}
	}

	return nil
}

// 头反转
func change_heads(file_name string) error {
	decode_file, _, newfile := openfile(file_name)
	err, newfile := exchange(decode_file, 0, 1, 2, 1, 8, file_name, newfile)
	if err != nil {
		return err
	}
	err, newfile = exchange(decode_file, 1, 1, 3, 1, 8, file_name, newfile)
	if err != nil {

		return err
	}
	generates_picture(file_name, newfile)
	return nil
}
