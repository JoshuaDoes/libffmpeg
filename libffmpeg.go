package Libffmpeg

import (
	"fmt"
	"strconv"
	"time"

	"github.com/JoshuaDoes/crunchio"
	"github.com/JoshuaDoes/ffmpeg"
)

var (
	filters   []*ffmpeg.Filter
	instances []*ffmpeg.Ffmpeg
	libraries []*Libffmpeg
)

type Libffmpeg struct {
	i       int
	in, out *crunchio.Buffer
	done    bool
}

func NewLibffmpeg(codec, format string) *Libffmpeg {
	if filters == nil {
		filters = make([]*ffmpeg.Filter, 0)
	}
	if instances == nil {
		instances = make([]*ffmpeg.Ffmpeg, 0)
	}
	if libraries == nil {
		libraries = make([]*Libffmpeg, 0)
	}
	i := len(instances)

	ff := ffmpeg.NewFFmpeg(codec, format)
	ff.SetName(fmt.Sprintf("%d", i))
	ff.SetOnExit(onExit)
	instances = append(instances, ff)

	in := ff.GetBufferAudioIn()
	out := ff.GetBufferAudioOut()

	lf := &Libffmpeg{
		i:   i,
		in:  in,
		out: out,
	}
	libraries = append(libraries, lf)

	return lf
}

func (lf *Libffmpeg) get() *ffmpeg.Ffmpeg {
	return instances[lf.i]
}

func (lf *Libffmpeg) IsDone() bool {
	return lf.done
}

func onExit(ff *ffmpeg.Ffmpeg) {
	i, err := strconv.Atoi(ff.GetName())
	if err != nil {
		panic(err)
	}
	lf := libraries[i]
	lf.done = true
}

func (lf *Libffmpeg) AddArgIn(arg string) {
	lf.get().AddArgsIn(arg)
}

func (lf *Libffmpeg) AddArgOut(arg string) {
	lf.get().AddArgsOut(arg)
}

func (lf *Libffmpeg) GetFFmpeg() string {
	return lf.get().GetFFmpeg()
}

func (lf *Libffmpeg) SetFFmpeg(path string) {
	lf.get().SetFFmpeg(path)
}

func (lf *Libffmpeg) GetLibraryPath() string {
	return lf.get().GetLibraryPath()
}

func (lf *Libffmpeg) SetLibraryPath(path string) {
	lf.get().SetLibraryPath(path)
}

func (lf *Libffmpeg) Start() {
	if err := lf.get().Start(); err != nil {
		panic(err)
	}
}

func (lf *Libffmpeg) Close() {
	if err := lf.get().Close(); err != nil {
		panic(err)
	}
}

func (lf *Libffmpeg) IsRunning() bool {
	return lf.get().IsRunning()
}

func (lf *Libffmpeg) Run() {
	if err := lf.get().Run(); err != nil {
		panic(err)
	}
}

func (lf *Libffmpeg) GetStats() string {
	return lf.get().GetBufferStats().String()
}

func (lf *Libffmpeg) SetBufferLength(milli int) {
	lf.get().SetBufferLength(time.Millisecond * time.Duration(milli))
}

func (lf *Libffmpeg) SetBufferSize(n int64) {
	lf.get().SetBufferSize(n)
}

func (lf *Libffmpeg) SetInput(input string) {
	lf.get().SetInput(input)
	if input != "" {
		lf.in = nil
	} else {
		lf.in = lf.get().GetBufferAudioIn()
	}
}

func (lf *Libffmpeg) SetOutput(output string) {
	lf.get().SetOutput(output)
	if output != "" {
		lf.out = nil
	} else {
		lf.out = lf.get().GetBufferAudioOut()
	}
}

func (lf *Libffmpeg) SetInputCodec(codec string) {
	lf.get().SetInputCodec(codec)
}

func (lf *Libffmpeg) SetOutputCodec(codec string) {
	lf.get().SetOutputCodec(codec)
}

func (lf *Libffmpeg) SetInputFormat(format string) {
	lf.get().SetInputFormat(format)
}

func (lf *Libffmpeg) SetOutputFormat(format string) {
	lf.get().SetOutputFormat(format)
}

func (lf *Libffmpeg) SetInputRate(rate int) {
	lf.get().SetInputRate(rate)
}

func (lf *Libffmpeg) SetOutputRate(rate int) {
	lf.get().SetOutputRate(rate)
}

func (lf *Libffmpeg) SetInputBitrate(bitrate int) {
	lf.get().SetInputBitrate(bitrate)
}

func (lf *Libffmpeg) SetOutputBitrate(bitrate int) {
	lf.get().SetOutputBitrate(bitrate)
}

func (lf *Libffmpeg) SetThreads(threads int) {
	lf.get().SetThreads(threads)
}

func (lf *Libffmpeg) SetPrecision(precision string) {
	lf.get().SetPrecision(precision)
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

func (lf *Libffmpeg) CommandRaw() string {
	return lf.get().String()
}

func (lf *Libffmpeg) Error() string {
	if err := lf.get().Error(); err != nil {
		return fmt.Sprintf("%v", err)
	}
	return ""
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

func (lf *Libffmpeg) Read(max int) []byte {
	if max < 0 {
		lf.out.Close()
		return nil
	}

	bytes := make([]byte, max)
	read, err := lf.out.Read(bytes)
	if err != nil {
		bytes = nil
	} else {
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
