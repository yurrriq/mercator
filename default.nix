{ stdenv, buildGoPackage, fetchFromGitHub }:

buildGoPackage rec {
  name = "mercator-unstable-${version}";
  version = "2018-04-13";
  rev = "d6c8de37b24801d1268284859cbd748f58320c91";

  goPackagePath = "github.com/shanesiebken/mercator";

  src = fetchFromGitHub {
    inherit rev;
    owner = "shanesiebken";
    repo = "mercator";
    sha256 = "1sl9f0qk25pzhbhkv5rlasgsdg025p5pxm03ralinznpmgidlr22";
  };

  goDeps = ./deps.nix;

  meta = with stdenv.lib; {
    inherit (src.meta) homepage;
    description = "A templating wrapper for Helm charts, introducing the concept of Chart \"Projections\"";
    platforms = platforms.all;
    license = license.asl20;
  };
}
