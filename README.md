你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# photosort
A small program written in Golang to copy and sort photos by year and month for a source to a dest folder.

This is my first attempt at the go language. Feel free to improve it through pull requests!

***Important: This has only been tested on Windows 10, let me know how it behaves in other OS.***

### Functionality
- Recursively visits all files and folders in the srcFolder directoy tree
- Analyses if files have Exif date information
- If so it copies the media file to the corresponding folder in the destfolder, creating the necessary folders such as Year and Month
- If the file is a media file but has no information regarding the capture date, it copies the file into a folder called Others in the destination folder
- It ignores any duplicated file or non media files

### Media files supported
```.tiff .tif .gif .jpeg .jpg .png .raw .webm .mkv .avi .mov .wmv .mp4 .m4v .mpg .mp2 .mpeg```

### Build
```
go get github.com/rwcarlsen/goexif/exif
go build photosort.go
```

### Run
```
./photosort [srcFolder] [destFolder]
```
