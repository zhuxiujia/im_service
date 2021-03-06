/**
 * Copyright (c) 2014-2015, GoBelieve
 * All rights reserved.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */

package ims

import "os"
import "fmt"
import "bytes"
import "flag"
import "encoding/binary"
import "path/filepath"
import "strings"
import "strconv"
import log "github.com/golang/glog"

var root string

func init() {
	flag.StringVar(&root, "root", "", "root")
}

func checkRoot(root string) {
	pattern := fmt.Sprintf("%s/message_*", root)
	files, _ := filepath.Glob(pattern)
	block_NO := 0 //begin from 0
	for _, f := range files {
		base := filepath.Base(f)
		if strings.HasPrefix(base, "message_") {
			if !CheckFile(f) {
				log.Warning("check file failure")
				r := truncateFile(f)
				log.Info("truncate file:", r)
			} else {
				log.Infof("check file pass:%s", f)
			}
			b, err := strconv.ParseInt(base[8:], 10, 64)
			if err != nil {
				log.Fatal("invalid message file:", f)
			}

			if int(b) > block_NO {
				block_NO = int(b)
			}
		}
	}

}

func truncateFile(file_path string) bool {
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal("open file:", err)
	}

	file_size, err := file.Seek(0, os.SEEK_END)
	if err != nil {
		log.Fatal("seek file")
	}

	if file_size == HEADER_SIZE {
		return true
	}

	if file_size < HEADER_SIZE {
		return false
	}

	offset := int64(4)

	for {
		_, err = file.Seek(file_size-offset, os.SEEK_SET)
		if err != nil {
			log.Fatal("seek file")
		}

		mf := make([]byte, 4)
		n, err := file.Read(mf)
		if err != nil || n != 4 {
			log.Fatal("read file err:", err)
		}
		buffer := bytes.NewBuffer(mf)
		var m int32
		binary.Read(buffer, binary.BigEndian, &m)

		if int(m) == MAGIC {
			log.Infof("file name:%s size:%d truncated:%d passed", file_path, file_size, file_size-offset+4)
			return true
		}

		offset += 4
	}

	return false

}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func RunTruncate() {
	flag.Parse()

	if len(root) == 0 {
		log.Info("trunncate imsroot")
		return
	}

	if !IsDir(root) {
		log.Info(root, "is not dir")
		return
	}
	log.Info("checking root:", root)
	checkRoot(root)
}
