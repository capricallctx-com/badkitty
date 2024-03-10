Name:           badkitty
Version:        1.0+MANUALBUILD
Release:        1%{?dist}
Summary:        Hello World example implemented in C

License:        MIT
URL:            https://github.com/naveenrajm7/%{name}
Source0:        https://github.com/naveenrajm7/%{name}/%{name}-%{version}.tar.gz

BuildRequires:  go


%description
BadKitty is a web server for modern stacks.

%prep
%setup -q

%build
go build -o badkitty


%install
%make_install


%files
%license LICENSE
%{_bindir}/%{name}
