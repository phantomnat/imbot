swc

package: `com.com2us.chronicles.android.google.us.normal`

find current activity
```
adb shell
dumpsys window | grep -E 'mCurrentFocus|mFocusedApp'
```

stop app

```shell
am force-stop com.com2us.chronicles.android.google.us.normal
```

monkey -p com.com2us.chronicles.android.google.us.normal


cmd package resolve-activity -c android.intent.category.LAUNCHER com.com2us.chronicles.android.google.us.normal |sed -n '/name=/s/^.*name=//p'


com.com2us.SMMO.ChroniclesActivity

PACKAGE="com.com2us.chronicles.android.google.us.normal"
ACTIVITY="com.com2us.SMMO.ChroniclesActivity"
am start -n "${PACKAGE}/${ACTIVITY}"

PACKAGE="com.com2us.chronicles.android.google.us.normal"
ACTIVITY="com.com2us.SMMO.ChroniclesActivity"
am start -n "${PACKAGE}/${ACTIVITY}"
adb shell 