﻿trap {
  write-error $_
  exit 1
}

$env:GOPATH = Join-Path -Path $PWD "bosh-dns-release"
$env:PATH = $env:GOPATH + "/bin;" + $env:PATH

cd $env:GOPATH
cd $env:GOPATH/src/bosh-dns

Push-Location "acceptance_tests\dns-acceptance-release\src\test-recursor"
$env:TEST_RECURSOR_BINARY = Join-Path -Path $PWD -ChildPath "test-recursor.exe"
go.exe build -o $env:TEST_RECURSOR_BINARY .
Pop-Location

go.exe run github.com/onsi/ginkgo/ginkgo -p -r -race -keepGoing -randomizeAllSpecs -randomizeSuites dns healthcheck

if ($LastExitCode -ne 0)
{
    Write-Error $_
    exit 1
}

go.exe run github.com/onsi/ginkgo/ginkgo -r -race -keepGoing -randomizeAllSpecs -randomizeSuites integration-tests

if ($LastExitCode -ne 0)
{
    Write-Error $_
    exit 1
}
