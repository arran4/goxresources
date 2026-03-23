Name:           xresources
Version:        %{?version}%{!?version:0.0.0}
Release:        1%{?dist}
Summary:        xresources summary
License:        MIT
Source0:        %{name}-%{version}.tar.gz

%description
xresources description.

%prep
%autosetup

%build
# build steps here

%install
mkdir -p %{buildroot}/usr/bin

%files
/usr/bin/*

%changelog
* Thu Mar 04 2026 CI Bot <ci@example.com> - %{version}-1
- Automated source build
