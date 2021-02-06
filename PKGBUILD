# Maintainer: Alexander Brevig <alexanderbrevig@gmail.com>
pkgname=marauder
pkgver=2021.1
pkgrel=1
pkgdesc='Documentation tool for terminal commands'
arch=('x86_64')
url="https://alexanderbrevig.com/$pkgname"
license=('MIT')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/AlexanderBrevig/marauder/archive/v$pkgver.tar.gz")

prepare(){
  cd "$srcdir/$pkgname-$pkgver"
  mkdir -p build/
}

build() {
  cd "$srcdir/$pkgname-$pkgver"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
  go build -o build ./cmd/...
}

check() {
  cd "$srcdir/$pkgname-$pkgver"
  #go test ./...
}

package() {
  cd "$srcdir/$pkgname-$pkgver"
  install -Dm755 build/$pkgname "$pkgdir"/usr/bin/$pkgname
}
md5sums=('4262ca3c1983168add96d8193c0454f8')
md5sums=('29b9625e78c6a3a4a1b01ed1a9575085')
