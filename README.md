# ncs-pkg-bumper

Simple package version bumper written in [Go Programming Language](https://go.dev/) for [Cisco NSO](https://developer.cisco.com/docs/nso/) packages.

## Buid

```bash
go build
```

## Usage

```bash
ncs-pkg-bumper -p test-package/package-meta-data.xml -m minor
```

**OR**

```bash
# looks for 'package-meta-data.xml' in specified directory
ncs-pkg-bumper -p test-package -m minor
```

**OR**

```bash
# looks for 'package-meta-data.xml' in current directory
ncs-pkg-bumper -m minor
```

## Example

**Command**

```bash
ncs-pkg-bumper -p test-package/package-meta-data.xml -m minor
```

**Console output**

```text
--- Package: test-package ---
Current version: 1.0.0
New version: 1.1.0
```

**Diff**

```xml
<ncs-package xmlns="http://tail-f.com/ns/ncs-packages">
  <name>test-package</name>
- <package-version>1.0.0</package-version>
+ <package-version>1.1.0</package-version>
  <description>ncs-pkg-bumper test package</description>
  <ncs-min-version>6.1</ncs-min-version>

  <!-- Some comment -->
  <!-- same package, multiple services, data providers etc -->

  <component>
    <name>MyComponent1</name>
    <callback>
      <java-class-name>com.example.test.testClass</java-class-name>
    </callback>
  </component>
  <component>
    <name>test</name>
    <application>
      <python-class-name>test.Main</python-class-name>
    </application>
  </component>
</ncs-package>
```
