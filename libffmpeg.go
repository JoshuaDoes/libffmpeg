package Libffmpeg

import (
	"fmt"
	"os"
	"time"

	"github.com/JoshuaDoes/crunchio"
	"github.com/JoshuaDoes/ffmpeg"
)

var (
	filters   []*ffmpeg.Filter
	instances []*ffmpeg.Ffmpeg
)

type Libffmpeg struct {
	i       int
	in, out *crunchio.Buffer
}

func NewLibffmpeg(codec, format string) *Libffmpeg {
	if filters == nil {
		filters = make([]*ffmpeg.Filter, 0)
	}
	if instances == nil {
		instances = make([]*ffmpeg.Ffmpeg, 0)
	}

	ff := ffmpeg.NewFFmpeg(codec, format)
	instances = append(instances, ff)

	return &Libffmpeg{
		i:   len(instances) - 1,
		in:  ff.GetBufferAudioIn(),
		out: ff.GetBufferAudioOut(),
	}
}

func (lf *Libffmpeg) get() *ffmpeg.Ffmpeg {
	return instances[lf.i]
}

func (lf *Libffmpeg) GetWD() string {
	wd, _ := os.Getwd()
	return wd
}

func (lf *Libffmpeg) GetFFmpeg() string {
	return lf.get().GetFFmpeg()
}

func (lf *Libffmpeg) SetFFmpeg(path string) {
	lf.get().SetFFmpeg(path)
}

func (lf *Libffmpeg) Start() {
	if err := lf.get().Start(); err != nil {
		panic(err)
	}
}

func (lf *Libffmpeg) Run() {
	if err := lf.get().Run(); err != nil {
		panic(err)
	}
}

func (lf *Libffmpeg) Read(max int) []byte {
	fmt.Println("FFmpeg: 1")
	if max < 0 {
		lf.out.Close()
		return nil
	}
	fmt.Println("FFmpeg: 2")

	bytes := make([]byte, max)
	read, err := lf.out.Read(bytes)
	if err != nil {
		fmt.Println("FFmpeg: 3", err)
		bytes = nil
	} else {
		fmt.Println("FFmpeg: 4", read)
		bytes = bytes[:read]
	}
	return bytes
}

func (lf *Libffmpeg) Write(src []byte) int {
	if src == nil {
		lf.in.Close()
		return -1
	}

	wrote, err := lf.in.Write(src)
	if err != nil {
		return -1
	}
	return wrote
}

func (lf *Libffmpeg) SetMetadata(key, value string) {
	lf.get().SetMetadata(key, value)
}

func (lf *Libffmpeg) SetInputChannels(channels int) {
	lf.get().SetInputChannels(channels)
}

func (lf *Libffmpeg) SetOutputChannels(channels int) {
	lf.get().SetOutputChannels(channels)
}

func (lf *Libffmpeg) SetOutputBitrate(bitrate int) {
	lf.get().SetOutputBitrate(bitrate)
}

func (lf *Libffmpeg) SetOutputRate(rate int) {
	lf.get().SetOutputRate(rate)
}

func (lf *Libffmpeg) SetPrecision(precision string) {
	lf.get().SetPrecision(precision)
}

func (lf *Libffmpeg) SetThreads(threads int) {
	lf.get().SetThreads(threads)
}

func (lf *Libffmpeg) SetBufferLength(milli int) {
	lf.get().SetBufferLength(time.Duration(milli))
}

func (lf *Libffmpeg) SetBufferSize(size int) {
	lf.get().SetBufferSize(int64(size))
}

func (lf *Libffmpeg) NewFilter() int {
	filters = append(filters, lf.get().NewFilter())
	return int(len(filters) - 1)
}

func (lf *Libffmpeg) FilterSetOutputChannelsMono(f, channel int) {
	filters[f].SetOutputChannels(int(channel))
}

func (lf *Libffmpeg) FilterSetOutputChannelsStereo(f, left, right int) {
	filters[f].SetOutputChannels(int(left), int(right))
}

func (lf *Libffmpeg) FilterAddHighpass(f, freq int) {
	filters[f].Add(ffmpeg.NewFilterHighpass(int(freq)))
}

func (lf *Libffmpeg) FilterAddLowpass(f, freq int) {
	filters[f].Add(ffmpeg.NewFilterLowpass(int(freq)))
}

func (lf *Libffmpeg) FilterAddEqualizer(f, freq, width int, gain float64) {
	filters[f].Add(ffmpeg.NewFilterEqualizer(freq, width, gain))
}

func (lf *Libffmpeg) FilterAddVolumeGain(f int, gain float64) {
	filters[f].Add(ffmpeg.NewFilterVolumeGain(gain))
}
