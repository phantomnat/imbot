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

###  ResourceType.Raw
| resource name | qty |
| --- | --- |
|   |   |
| branch | 18950 |
| leaf | 1800 |
| sturdy tree trunk | 10650 |
| tough leaf | 1650 |
|   |   |
| adamantium mineral | 6000 |
| emerald ore | 2700 |
| fluorite mineral | 300 |
| gold mineral | 6000 |
| obsidian mineral | 300 |
| peridot ore | 600 |
| sapphire ore | 14850 |
| tritium mineral | 1050 |
|   |   |
| shoerekly's filament | 5 |
| tallades carapace | 5 |

###  ResourceType.Processing
| resource name | qty |
| --- | --- |
| adamantium ingot | 2000 |
| emerald | 900 |
| fluorite ingot | 100 |
| gold ingot | 2000 |
| magic_craft_tool | 1050 |
| normal lumber | 900 |
| normal_string | 600 |
| obsidian | 100 |
| peridot | 200 |
| premium hardwork craft tool | 300 |
| sapphire | 4950 |
| solid lumber | 3250 |
| tritium | 350 |

###  ResourceType.Alchemy
| resource name | qty |
| --- | --- |
| spell_book_iv_atk | 20 |
| spell_book_iv_util | 20 |
| supreme green handwork gem | 20 |
| supreme moon stone (1) | 10 |
| supreme moon stone (2) | 5 |
| supreme moon stone (3) | 5 |

Process finished with exit code 0
