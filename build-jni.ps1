# Generate Android App Resource bindings for Android 12 (requires install of Android SDK 31 or newer)
go get -u golang.org/x/mobile/bind@latest
gomobile bind -o libffmpeg.aar -target=android -androidapi 31 .
