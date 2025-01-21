# Generate Android App Resource bindings for Android 15 (requires install of Android SDK 35)
go get -u golang.org/x/mobile/bind@latest
gomobile bind -o libffmpeg.aar -target=android -androidapi 35 .
