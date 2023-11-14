import dataclasses
from enum import Enum, auto
from typing import Dict


class ResourceType(Enum):
    Raw = auto()
    Shopping = auto()
    Processing = auto()
    Alchemy = auto()
    Cooking = auto()



@dataclasses.dataclass
class Resource:
    name: str
    type: ResourceType
    sub_type: str = None
    qty: int = 0
    deps: Dict[str, int] = None


resources = {}
# raw

# drop
magic_core = 'magic core'
faint_magic_core = 'faint magic core'
powerful_magic_core = 'powerful magic core'
flame_heart = 'flame heart'
earth_breath = 'earth breath'
dawn_frost = 'dawn frost'
chicken = 'chicken'
beef = 'beef'
lamb = 'lamb'
meat = 'meat'
tallades_carapace = 'tallades carapace'
shoereklys_filament = "shoerekly's filament"
fire_spirit_heart = 'fire spirit heart'
naraka_s_claw = "naraka's claw"
extreme_frost = 'extreme frost'
vita_tau_s_red_gem = "vita-tau's red gem"
raw_drop_items = [
    # drop
    chicken,
    beef,
    lamb,
    meat,

    magic_core,
    faint_magic_core,
    powerful_magic_core,
    flame_heart,
    earth_breath,
    dawn_frost,
    tallades_carapace,
    shoereklys_filament,
    fire_spirit_heart,
    naraka_s_claw,
    extreme_frost,
    vita_tau_s_red_gem,
]

# gathering
branch = 'branch'
sturdy_tree_trunk = 'sturdy tree trunk'
tough_leaf = 'tough leaf'
leaf = 'leaf'
sage_leaf = 'sage_leaf'
saffron_petal = 'saffron petal'
eternal_leaf = 'eternal_leaf'
four_leaf_clover = '4-leaf clover'
chamomile_bud = 'chamomile bud'
ice_mango = 'ice_mango'
apple = 'apple'
sweet_apple = 'sweet apple'
strawberry = 'strawberry'
lavender_stem = 'lavender_stem'
mango = 'mango'
basil_leaf = 'basil leaf'
kiwi = 'kiwi'
raw_gathering_items = [
    # gathering
    branch,
    sturdy_tree_trunk,
    leaf,
    tough_leaf,
    sage_leaf,
    saffron_petal,
    eternal_leaf,
    four_leaf_clover,
    chamomile_bud,
    ice_mango,
    apple,
    sweet_apple,
    strawberry,
    lavender_stem,
    mango,
    basil_leaf,
    kiwi,
]

# mining
unmelted_ice = 'unmelted ice'
sapphire_ore = 'sapphire ore'
aquamarine_ore = 'aquamarine ore'
emerald_ore = 'emerald ore'
peridot_ore = 'peridot ore'
gold_mineral = 'gold mineral'
fluorite_mineral = 'fluorite mineral'
obsidian_mineral = 'obsidian mineral'
adamantium_mineral = 'adamantium mineral'
tritium_mineral = 'tritium mineral'
raw_mining_items = [
    # mining
    unmelted_ice,
    sapphire_ore,
    aquamarine_ore,
    emerald_ore,
    peridot_ore,
    gold_mineral,
    fluorite_mineral,
    obsidian_mineral,
    adamantium_mineral,
    tritium_mineral,
]

# fishing
crab = 'crab'
blowfish = 'blowfish'
swallowtail_dart = 'swallowtail dart'
calm_mackerel = 'calm mackerel'
mandarin_fish = 'mandarin fish'
barracuda_fish = 'barracuda fish'
catfish = 'catfish'
raw_fishing_items = [
    crab,
    blowfish,
    swallowtail_dart,
    calm_mackerel,
    mandarin_fish,
    barracuda_fish,
    catfish,
]
gold = 'gold'
raw_currencies = [
    'gold'
]
raw_items_with_sub_types = [
    (raw_currencies, 'currency'),
    (raw_gathering_items, 'gathering'),
    (raw_mining_items, 'mining'),
    (raw_fishing_items, 'fishing'),
    (raw_drop_items, 'drop'),
]
raw_item_count = 1
for items, sub_type in raw_items_with_sub_types:
    for name in items:
        resources[name] = Resource(name=name, type=ResourceType.Raw, sub_type=f'{raw_item_count}_{sub_type}')
    raw_item_count += 1

# shopping
olive_oil = 'olive oil'
flour = 'flour'
sugar = 'sugar'
salt = 'salt'
pepper = 'pepper'
natural_water = 'natural water'
shopping_items = [
    [natural_water, {gold: 80}],
    [pepper, {gold: 1000}],
    [olive_oil, {gold: 1000}],
    [flour, {gold: 1100}],
    [sugar, {gold: 100}],
    [salt, {gold: 100}],
]

# processing

aquamarine = 'aquamarine'
sapphire = 'sapphire'
emerald = 'emerald'
peridot = 'peridot'
obsidian = 'obsidian'
fluorite_ingot = 'fluorite ingot'
adamantium_ingot = 'adamantium ingot'
gold_ingot = 'gold ingot'
tritium = 'tritium'

magic_craft_tool = 'magic_craft_tool'
premium_handwork_craft_tool = 'premium hardwork craft tool'
normal_lumber = 'normal lumber'
solid_lumber = 'solid lumber'
normal_string = 'normal_string'

processing_items = [
    [gold_ingot, {gold: 1000, gold_mineral: 3}],
    [aquamarine, {gold: 1000, aquamarine_ore: 3}],

    [sapphire, {gold: 1500, sapphire_ore: 3}],
    [emerald, {gold: 1500, emerald_ore: 3}],
    [peridot, {gold: 3000, peridot_ore: 3}],

    [adamantium_ingot, {gold: 1500, adamantium_mineral: 3}],
    [tritium, {gold: 5000, tritium_mineral: 3}],
    [obsidian, {gold: 5000, obsidian_mineral: 3}],
    [fluorite_ingot, {gold: 4000, fluorite_mineral: 3}],

    [normal_lumber, {gold: 1500, sturdy_tree_trunk: 1, branch: 3}],
    [solid_lumber, {gold: 5000, sturdy_tree_trunk: 3, branch: 5}],
    [normal_string, {gold: 1500, tough_leaf: 1, leaf: 3}],
    [premium_handwork_craft_tool, {gold: 5000, emerald: 3, normal_string: 2, normal_lumber: 3}],
    # lv 7
    [magic_craft_tool, {gold: 7500, sapphire: 3, tough_leaf: 1, solid_lumber: 3}],
]

# cooking
cooking_sub_type_ingredient = 'ingredient'
cooking_sub_type_food = 'food'

fried_meat = 'fried meat'
chamomile_pho = 'chamomile pho'
mango_sorbet = 'mango sorbet'
premium_swallowtail_dart_dish = 'premium swallowtail dart dish'
strawberry_pie = 'strawberry pie'
grilled_calm_mackerel = 'grilled calm mackerel'
crab_gimbap = 'crab gimbap'
mandarin_fish_soup = 'mandarin fish soup'
kiwi_salad = 'kiwi salad'
apple_bread = 'apple bread'
barracuda_dish = 'barracuda dish'
grilled_catfish_skewer = 'grilled catfish skewer'
apple_fries = 'apple fries'
ham_sandwich = 'ham sandwich'
lavender_tea = 'lavender tea'
basil_tea = 'basil tea'

premium_flour = 'premium flour'
premium_sugar = 'premium sugar'
premium_salt = 'premium salt'
finest_flour = 'finest flour'
finest_sugar = 'finest sugar'
finest_salt = 'finest salt'

premium_all_purpose_ingredient = 'premium all-purpose ingredient'
# premium_all_purpose_ingredient_1 = 'premium all-purpose ingredient (1)'
# premium_all_purpose_ingredient_2 = 'premium all-purpose ingredient (2)'
# premium_all_purpose_ingredient_3 = 'premium all-purpose ingredient (3)'

cooking_items = [
    [premium_flour, cooking_sub_type_ingredient, {gold: 250, lavender_stem: 5, flour: 3}],
    [premium_sugar, cooking_sub_type_ingredient, {gold: 250, mango: 5, sugar: 3}],
    [premium_salt, cooking_sub_type_ingredient, {gold: 250, basil_leaf: 5, salt: 3}],

    [finest_flour, cooking_sub_type_ingredient, {gold: 500, saffron_petal: 5, eternal_leaf: 1, flour: 5}],
    [finest_sugar, cooking_sub_type_ingredient, {gold: 500, apple: 5, sweet_apple: 5, sugar: 5}],
    [finest_salt, cooking_sub_type_ingredient, {gold: 500, chamomile_bud: 5, four_leaf_clover: 1, salt: 5}],

    [premium_all_purpose_ingredient, cooking_sub_type_ingredient, {sage_leaf: 3, beef: 1, blowfish: 1}],

    [fried_meat, cooking_sub_type_food, {gold: 10_000, beef: 8, lamb: 8, finest_flour: 3, olive_oil: 5}],
    [chamomile_pho, cooking_sub_type_food, {gold: 10_000, chamomile_bud: 10, four_leaf_clover: 3, premium_all_purpose_ingredient: 2, finest_flour: 2}],
    [mango_sorbet, cooking_sub_type_food, {gold: 10_000, ice_mango: 8, unmelted_ice: 5, finest_sugar: 3}],
    [premium_swallowtail_dart_dish, cooking_sub_type_food, {gold: 10_000, sage_leaf: 6, swallowtail_dart: 4, premium_all_purpose_ingredient: 2, finest_salt: 2}],
    [strawberry_pie, cooking_sub_type_food, {gold: 5_000, strawberry: 12, premium_flour: 5, premium_sugar: 3}],
    [grilled_calm_mackerel, cooking_sub_type_food, {gold: 500, calm_mackerel: 5, premium_salt: 3}],
    [crab_gimbap, cooking_sub_type_food, {gold: 10_000, crab: 6, chicken: 5, premium_all_purpose_ingredient: 2, finest_flour: 3}],

    [mandarin_fish_soup, cooking_sub_type_food, {gold: 1_000, salt: 1, mandarin_fish: 2}],
    [kiwi_salad, cooking_sub_type_food, {gold: 2_000, kiwi: 5, sage_leaf: 3, olive_oil: 1}],
    [apple_bread, cooking_sub_type_food, {gold: 2_000, apple: 5, flour: 3, sugar: 3}],
    [barracuda_dish, cooking_sub_type_food, {gold: 1_500, barracuda_fish: 4, salt: 1}],
    [grilled_catfish_skewer, cooking_sub_type_food, {gold: 2000, catfish:5, salt: 2}],
    [apple_fries, cooking_sub_type_food, {gold: 2000, apple: 8, flour: 1}],
    [ham_sandwich, cooking_sub_type_food, {gold: 1000, meat: 3, flour: 1}],
    [lavender_tea, cooking_sub_type_food, {gold: 1000, lavender_stem: 3, natural_water: 1}],
    [basil_tea, cooking_sub_type_food, {gold: 1000, basil_leaf: 6, sugar: 1, natural_water: 1}],
]

# alchemy

spell_book_i_atk = 'spell book i atk'
spell_book_i_def = 'spell book i def'
spell_book_i_util = 'spell book i util'

spell_book_ii_atk = 'spell book ii atk'
spell_book_ii_def = 'spell book ii def'
spell_book_ii_util = 'spell book ii util'

spell_book_iii_atk = 'spell book iii atk'
spell_book_iii_def = 'spell book iii def'
spell_book_iii_util = 'spell book iii util'

spell_book_iv_atk = 'spell_book_iv_atk'
spell_book_iv_def = 'spell_book_iv_def'
spell_book_iv_util = 'spell_book_iv_util'

supreme_blue_handwork_gem = 'supreme blue handwork gem'
supreme_green_handwork_gem = 'supreme green handwork gem'

supreme_moon_stone_1 = 'supreme moon stone (1)'
supreme_moon_stone_2 = 'supreme moon stone (2)'
supreme_moon_stone_3 = 'supreme moon stone (3)'

alchemy_items = [
    [spell_book_i_atk, {faint_magic_core: 1, tough_leaf: 1, gold: 3000}],
    [spell_book_i_def, {faint_magic_core: 1, tough_leaf: 1, gold: 3000}],

    [spell_book_ii_atk, {faint_magic_core: 3, tough_leaf: 5, gold: 5000}],
    [spell_book_ii_def, {faint_magic_core: 3, tough_leaf: 5, gold: 5000}],
    [spell_book_ii_util, {faint_magic_core: 3, tough_leaf: 5, gold: 5000}],

    [spell_book_iii_atk, {magic_craft_tool: 3, magic_core: 5, faint_magic_core: 30, flame_heart: 1, gold: 12000}],
    [spell_book_iii_def, {magic_craft_tool: 3, magic_core: 5, faint_magic_core: 30, earth_breath: 1, gold: 12000}],
    [spell_book_iii_util, {magic_craft_tool: 3, magic_core: 5, faint_magic_core: 30, dawn_frost: 1, gold: 12000}],

    [spell_book_iv_atk, {gold: 20_000, magic_craft_tool: 25, powerful_magic_core: 5, magic_core: 25, fire_spirit_heart: 1, naraka_s_claw: 1 }],
    [spell_book_iv_def, {gold: 20_000, magic_craft_tool: 25, powerful_magic_core: 5, magic_core: 25}],
    [spell_book_iv_util, {gold: 20_000, magic_craft_tool: 25, powerful_magic_core: 5, magic_core: 25, extreme_frost: 1, vita_tau_s_red_gem: 1}],

    [supreme_blue_handwork_gem, {gold: 20_000, sapphire: 90, aquamarine: 30, premium_handwork_craft_tool: 15}],
    [supreme_green_handwork_gem, {gold: 20_000, emerald: 90, peridot: 10, premium_handwork_craft_tool: 15}],

    [supreme_moon_stone_1, {gold: 20_000, adamantium_ingot: 100, gold_ingot: 100, tritium: 25}],
    [supreme_moon_stone_2, {gold: 20_000, adamantium_ingot: 100, gold_ingot: 100, tallades_carapace: 1}],
    [supreme_moon_stone_3, {gold: 20_000, adamantium_ingot: 100, gold_ingot: 100, shoereklys_filament: 1}],

]

for name, deps in shopping_items:
    resources[name] = Resource(name=name, type=ResourceType.Shopping, deps=deps)
for name, sub_type, deps in cooking_items:
    resources[name] = Resource(name=name, type=ResourceType.Cooking, deps=deps, sub_type=sub_type)
for name, deps in processing_items:
    resources[name] = Resource(name=name, type=ResourceType.Processing, deps=deps)
for name, deps in alchemy_items:
    resources[name] = Resource(name=name, type=ResourceType.Alchemy, deps=deps)

required_resources = {}

def recursive_add(name: str, qty: int):
    if name not in required_resources:
        required_resources[name] = 0
    required_resources[name] += qty

    if name not in resources or resources[name].deps is None:
        return

    for dep_name, dep_qty_per_piece in resources[name].deps.items():
        recursive_add(dep_name, qty * dep_qty_per_piece)


# resources[spell_book_i_atk].qty = 225

# resources[spell_book_ii_atk].qty = 25
# resources[spell_book_ii_def].qty = 25
# resources[spell_book_ii_util].qty = 25

# lv 6 alchemy
# resources[spell_book_iii_atk].qty = 3
# resources[spell_book_iii_def].qty = 3
# resources[spell_book_iii_util].qty = 4
# resources[supreme_blue_handwork_gem].qty = 6

# master

# resources[fluorite_ingot].qty = 100
# resources[obsidian].qty = 100
# resources[tritium].qty = 100
#
# resources[magic_craft_tool].qty = 50
# resources[solid_lumber].qty = 100
#
# resources[spell_book_iv_atk].qty = 20
# resources[spell_book_iv_util].qty = 20
#
# resources[supreme_moon_stone_1].qty = 10
# resources[supreme_moon_stone_2].qty = 5
# resources[supreme_moon_stone_3].qty = 5
#
# resources[supreme_green_handwork_gem].qty = 20

# cooking
#
# resources[fried_meat].qty = 80
# resources[chamomile_pho].qty = 80
# resources[mango_sorbet].qty = 80
# resources[premium_swallowtail_dart_dish].qty = 80
# resources[strawberry_pie].qty = 80
# resources[crab_gimbap].qty = 80
# resources[grilled_calm_mackerel].qty = 80

# book
# cooking
# resources[mandarin_fish_soup].qty = 20
# resources[kiwi_salad].qty = 20
# resources[apple_bread].qty = 20
# resources[barracuda_dish].qty = 15
# resources[grilled_catfish_skewer].qty = 15
# resources[apple_fries].qty = 10
# resources[ham_sandwich].qty = 10
# resources[lavender_tea].qty = 10
# resources[basil_tea].qty = 10

for name, resource in resources.items():
    if resource.qty == 0:
        continue

    if resource.name not in required_resources:
        required_resources[resource.name] = 0

    required_resources[resource.name] += resource.qty

    if resource.deps is None:
        continue

    for dep_name, dep_qty_per_piece in resource.deps.items():
        # dep_qty = resource.qty*dep_qty_per_piece
        # recursive search
        recursive_add(dep_name, resource.qty * dep_qty_per_piece)

main_type = None
sub_type = None
for name in sorted(required_resources.keys(), key=lambda x: (resources[x].type.value, resources[x].sub_type, x)):
    if main_type is None or main_type != resources[name].type:
        print()
        # print('------')
        main_type = resources[name].type
        print('### ', main_type)
        # print('------')
        print(f'| resource name | qty |')
        print(f'| --- | --- |')
        sub_type = None
    if sub_type != resources[name].sub_type:
        # print()
        sub_type = resources[name].sub_type
        print('|   |   |')

    print(f'| {name} | {required_resources[name]:,} |')
