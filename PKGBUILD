# Maintainer: Lon60 37051290+Lon60@users.noreply.github.com
pkgname=sysfail
pkgver=1.0.0
pkgrel=1
pkgdesc="Linux System Failure Simulator - Realistic kernel panic and emergency shell prank tool"
arch=('x86_64' 'aarch64')
url="https://github.com/yourusername/sysfail"
license=('MIT')
depends=()
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
    cd "$pkgname-$pkgver"
    export CGO_ENABLED=0
    export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external"
    make build
}

package() {
    cd "$pkgname-$pkgver"
    make DESTDIR="$pkgdir" install

    # Install license
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"

    # Install documentation
    install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}