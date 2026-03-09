from pathlib import Path

from PIL import Image, ImageDraw


ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "internal" / "assets" / "generated"
OUT.mkdir(parents=True, exist_ok=True)


def rgba(hex_code, alpha=255):
    hex_code = hex_code.lstrip("#")
    return tuple(int(hex_code[i : i + 2], 16) for i in (0, 2, 4)) + (alpha,)


INK = rgba("172038")
SKY_TOP = rgba("4d65d9")
SKY_MID = rgba("7fc9ff")
SKY_LOW = rgba("f6c86f")
SUN = rgba("fff4c7")
MOUNTAIN_FAR = rgba("394a88")
MOUNTAIN_MID = rgba("2d365f", 210)
MOUNTAIN_NEAR = rgba("1e2032", 220)
DIRT_DARK = rgba("4a2d2a")
DIRT_MID = rgba("8d5a43")
DIRT_LIGHT = rgba("d59a64")
LEAF = rgba("5d8c4f")
WOOD_DARK = rgba("5c3628")
WOOD_LIGHT = rgba("b77545")
STONE_DARK = rgba("5d6479")
STONE_LIGHT = rgba("a8b6cf")
SKIN = rgba("f8c59b")
HAIR = rgba("362430")
SCARF = rgba("ef5d60")
CLOTH = rgba("54708a")
SHOE = rgba("2b2230")
GUN_DARK = rgba("2f3348")
GUN_LIGHT = rgba("7d8aa8")
FIRE_1 = rgba("ff8c42")
FIRE_2 = rgba("ffd166")
ICE_1 = rgba("7ce7ff")
ICE_2 = rgba("e7fdff")
THUNDER_1 = rgba("f7ea5c")
THUNDER_2 = rgba("fff9c8")
HEART = rgba("ff6978")


def save(name, image):
    image.save(OUT / name)


def outlined_rect(draw, box, fill, outline=INK):
    draw.rectangle(box, fill=fill)
    draw.rectangle(box, outline=outline)


def outlined_rect_no_bottom(draw, box, fill, outline=INK):
    x0, y0, x1, y1 = box
    draw.rectangle(box, fill=fill)
    draw.line((x0, y0, x1, y0), fill=outline)
    draw.line((x0, y0, x0, y1), fill=outline)
    draw.line((x1, y0, x1, y1), fill=outline)


def gradient(width, height, top, mid, low):
    img = Image.new("RGBA", (width, height))
    px = img.load()
    assert px is not None
    for y in range(height):
        t = y / (height - 1)
        if t < 0.55:
            local = t / 0.55
            c0, c1 = top, mid
        else:
            local = (t - 0.55) / 0.45
            c0, c1 = mid, low
        px_color = tuple(int(c0[i] + (c1[i] - c0[i]) * local) for i in range(3)) + (
            255,
        )
        for x in range(width):
            px[x, y] = px_color
    return img


def player_frame(index):
    img = Image.new("RGBA", (48, 48), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)

    run_legs = [
        ((17, 31, 21, 43), (25, 30, 29, 43), (12, 24, 17, 31), (28, 23, 33, 30)),
        ((16, 31, 20, 43), (26, 29, 30, 43), (13, 25, 18, 31), (28, 22, 33, 29)),
        ((15, 31, 19, 43), (27, 31, 31, 43), (14, 26, 19, 32), (27, 21, 32, 29)),
        ((17, 30, 21, 43), (25, 31, 29, 43), (12, 23, 17, 30), (28, 24, 33, 31)),
        ((18, 29, 22, 43), (24, 31, 28, 43), (13, 21, 18, 28), (29, 25, 34, 32)),
        ((17, 31, 21, 43), (25, 29, 29, 43), (12, 24, 17, 30), (28, 22, 33, 29)),
    ]

    if index < 6:
        left_leg, right_leg, left_arm, right_arm = run_legs[index]
        head_y = 10
        scarf_shift = index % 3
    elif index == 6:
        left_leg, right_leg = (17, 30, 21, 43), (25, 30, 29, 43)
        left_arm, right_arm = (13, 20, 18, 27), (29, 18, 34, 25)
        head_y = 8
        scarf_shift = 2
    else:
        left_leg, right_leg = (18, 32, 22, 43), (24, 31, 28, 43)
        left_arm, right_arm = (12, 24, 17, 31), (29, 22, 34, 29)
        head_y = 11
        scarf_shift = 1

    outlined_rect(d, (18, head_y, 29, head_y + 11), SKIN)
    d.rectangle((18, head_y, 29, head_y + 4), fill=HAIR)
    d.rectangle((17, head_y + 3, 19, head_y + 7), fill=HAIR)
    d.point((26, head_y + 6), fill=INK)
    d.line((20, head_y + 11, 27, head_y + 11), fill=INK)

    outlined_rect(d, (18, 22, 29, 33), CLOTH)
    d.rectangle((20, 23, 27, 26), fill=rgba("7890a7"))
    d.rectangle((18, 22, 24, 25), fill=SCARF)
    d.polygon(
        [
            (18, 24),
            (11 - scarf_shift, 22 + scarf_shift),
            (13 - scarf_shift, 27 + scarf_shift),
            (18, 27),
        ],
        fill=SCARF,
    )
    d.line((11 - scarf_shift, 22 + scarf_shift, 18, 24), fill=INK)
    d.line(
        (11 - scarf_shift, 22 + scarf_shift, 13 - scarf_shift, 27 + scarf_shift),
        fill=INK,
    )
    d.line((13 - scarf_shift, 27 + scarf_shift, 18, 27), fill=INK)

    outlined_rect(d, left_arm, SKIN)
    outlined_rect(d, right_arm, SKIN)

    left_hand = (left_arm[0] - 1, left_arm[3] - 3, left_arm[2], left_arm[3])
    right_hand = (
        right_arm[2] - 1,
        right_arm[3] - 4,
        right_arm[2] + 2,
        right_arm[3] - 1,
    )
    outlined_rect(d, left_hand, SKIN)
    outlined_rect(d, right_hand, SKIN)

    pistol = [
        (right_arm[2] + 1, right_arm[1] + 1),
        (right_arm[2] + 7, right_arm[1] + 1),
        (right_arm[2] + 7, right_arm[1] + 4),
        (right_arm[2] + 5, right_arm[1] + 4),
        (right_arm[2] + 5, right_arm[1] + 7),
        (right_arm[2] + 3, right_arm[1] + 8),
        (right_arm[2] + 2, right_arm[1] + 4),
        (right_arm[2] + 1, right_arm[1] + 4),
    ]
    d.polygon(pistol, fill=GUN_DARK, outline=INK)
    d.rectangle(
        (right_arm[2] + 3, right_arm[1] + 2, right_arm[2] + 5, right_arm[1] + 2),
        fill=GUN_LIGHT,
    )
    d.rectangle(
        (right_arm[2] + 7, right_arm[1] + 2, right_arm[2] + 8, right_arm[1] + 2),
        fill=FIRE_2,
    )

    outlined_rect_no_bottom(d, left_leg, rgba("4d6478"))
    outlined_rect_no_bottom(d, right_leg, rgba("4d6478"))

    left_shoe = (left_leg[0] - 1, left_leg[3] - 1, left_leg[2] + 3, left_leg[3] + 1)
    right_shoe = (
        right_leg[0] - 1,
        right_leg[3] - 1,
        right_leg[2] + 3,
        right_leg[3] + 1,
    )
    outlined_rect(d, left_shoe, SHOE)
    outlined_rect(d, right_shoe, SHOE)

    return img


def enemy_ground_frame(index):
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    body_y = 14 + (index % 2)
    outlined_rect(d, (5, body_y, 24, 24 + (index % 2)), ICE_1)
    d.rectangle((8, body_y + 2, 21, body_y + 5), fill=ICE_2)
    d.polygon([(5, body_y + 4), (2, body_y + 1), (8, body_y)], fill=ICE_2)
    d.polygon([(24, body_y + 6), (28, body_y + 4), (24, body_y + 1)], fill=ICE_2)
    d.point((18, body_y + 5), fill=INK)
    feet = [(7, 25, 10, 29), (14, 25, 17, 29), (20, 25, 23, 29)]
    shift = [-1, 0, 1][index]
    for i, foot in enumerate(feet):
        x0, y0, x1, y1 = foot
        outlined_rect(
            d, (x0 + shift * (i - 1), y0, x1 + shift * (i - 1), y1), rgba("a0b8d4")
        )
    return img


def enemy_turret_frame(index):
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    outlined_rect(d, (8, 14, 23, 28), STONE_DARK)
    d.rectangle((10, 16, 21, 19), fill=STONE_LIGHT)
    flame_height = [11, 14, 12][index]
    d.polygon(
        [
            (15, 15),
            (12, 15 - flame_height // 2),
            (16, 4),
            (20, 15 - flame_height // 2),
            (17, 15),
        ],
        fill=FIRE_1,
    )
    d.polygon([(15, 14), (14, 11), (16, 7), (18, 11), (17, 14)], fill=FIRE_2)
    d.point((13, 19), fill=INK)
    d.point((18, 19), fill=INK)
    return img


def enemy_flyer_frame(index):
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    wing_l = [6, 4, 7][index]
    wing_r = [25, 27, 24][index]
    d.polygon([(11, 14), (wing_l, 10), (10, 20)], fill=THUNDER_1)
    d.polygon([(20, 14), (wing_r, 10), (21, 20)], fill=THUNDER_1)
    d.ellipse((10, 10, 21, 21), fill=THUNDER_1, outline=INK)
    d.ellipse((13, 13, 18, 18), fill=THUNDER_2)
    d.point((15, 14), fill=INK)
    d.line((15, 20, 12, 25), fill=THUNDER_1)
    d.line((15, 20, 18, 25), fill=THUNDER_1)
    return img


def projectile_fire(player_owned=True):
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    outer = FIRE_1 if player_owned else rgba("c85b2f")
    inner = FIRE_2 if player_owned else rgba("f0bc6a")
    d.ellipse((2, 3, 13, 12), fill=outer, outline=INK)
    d.ellipse((5, 5, 10, 10), fill=inner)
    d.polygon([(2, 8), (0, 6), (2, 4)], fill=outer)
    return img


def projectile_ice(player_owned=True):
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    outer = ICE_1 if player_owned else rgba("58b5d3")
    inner = ICE_2 if player_owned else rgba("d7f7ff")
    d.polygon([(3, 8), (8, 2), (13, 8), (8, 13)], fill=outer, outline=INK)
    d.polygon([(6, 8), (8, 5), (10, 8), (8, 11)], fill=inner)
    return img


def projectile_thunder(player_owned=True):
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    outer = THUNDER_1 if player_owned else rgba("cabf46")
    inner = THUNDER_2 if player_owned else rgba("fff0a2")
    d.polygon(
        [(5, 1), (10, 1), (8, 6), (12, 6), (6, 15), (7, 9), (3, 9)],
        fill=outer,
        outline=INK,
    )
    d.polygon([(7, 3), (9, 3), (8, 6), (10, 6), (7, 11), (7, 8), (5, 8)], fill=inner)
    return img


def fire_fx():
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.pieslice((1, 4, 31, 30), start=315, end=95, fill=FIRE_1, outline=INK)
    d.pieslice((7, 8, 27, 25), start=318, end=92, fill=FIRE_2)
    return img


def ice_fx():
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    points = [
        (16, 2),
        (22, 10),
        (30, 16),
        (22, 22),
        (16, 30),
        (10, 22),
        (2, 16),
        (10, 10),
    ]
    d.polygon(points, fill=ICE_1, outline=INK)
    d.polygon(
        [(16, 7), (20, 12), (25, 16), (20, 20), (16, 25), (12, 20), (7, 16), (12, 12)],
        fill=ICE_2,
    )
    return img


def thunder_fx():
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon(
        [(8, 1), (20, 1), (16, 10), (24, 10), (10, 31), (13, 18), (5, 18)],
        fill=THUNDER_1,
        outline=INK,
    )
    d.polygon(
        [(11, 4), (18, 4), (15, 10), (21, 10), (12, 23), (13, 15), (9, 15)],
        fill=THUNDER_2,
    )
    return img


def tile_ground():
    img = Image.new("RGBA", (32, 32), DIRT_DARK)
    d = ImageDraw.Draw(img)
    d.rectangle((0, 0, 31, 8), fill=LEAF)
    d.rectangle((0, 9, 31, 31), fill=DIRT_MID)
    for x in range(0, 32, 4):
        d.line((x, 8, x + 2, 10), fill=DIRT_LIGHT)
    for x, y in [(4, 15), (15, 18), (24, 13), (11, 25), (28, 23)]:
        d.rectangle((x, y, x + 2, y + 2), fill=DIRT_LIGHT)
    d.rectangle((0, 0, 31, 31), outline=INK)
    return img


def tile_platform():
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    outlined_rect(d, (0, 8, 31, 23), STONE_DARK)
    d.rectangle((1, 9, 30, 12), fill=STONE_LIGHT)
    for x in range(4, 28, 8):
        d.line((x, 8, x, 23), fill=INK)
    d.line((3, 23, 8, 28), fill=STONE_DARK)
    d.line((28, 23, 23, 28), fill=STONE_DARK)
    return img


def tile_crate():
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    outlined_rect(d, (2, 2, 29, 29), WOOD_LIGHT)
    d.line((2, 2, 29, 29), fill=WOOD_DARK)
    d.line((29, 2, 2, 29), fill=WOOD_DARK)
    d.rectangle((6, 6, 25, 10), fill=WOOD_DARK)
    d.rectangle((6, 21, 25, 25), fill=WOOD_DARK)
    return img


def tile_spike():
    img = Image.new("RGBA", (32, 32), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    spikes = [
        [(1, 31), (6, 12), (11, 31)],
        [(10, 31), (16, 7), (22, 31)],
        [(20, 31), (26, 12), (31, 31)],
    ]
    for spike in spikes:
        d.polygon(spike, fill=rgba("e3c2ff"), outline=INK)
    d.rectangle((0, 28, 31, 31), fill=STONE_DARK, outline=INK)
    return img


def heart_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon(
        [
            (8, 13),
            (3, 8),
            (3, 4),
            (5, 2),
            (7, 2),
            (8, 4),
            (9, 2),
            (11, 2),
            (13, 4),
            (13, 8),
        ],
        fill=HEART,
    )
    d.polygon([(4, 5), (5, 4), (6, 4)], fill=rgba("ffffff", 150))
    return img


def fire_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon([(8, 2), (4, 9), (5, 13), (8, 14), (11, 13), (12, 9)], fill=FIRE_1)
    d.polygon([(8, 7), (6, 10), (8, 12), (10, 10)], fill=FIRE_2)
    return img


def ice_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon([(8, 2), (4, 8), (8, 14), (12, 8)], fill=ICE_1)
    d.polygon([(8, 5), (6, 8), (8, 11), (10, 8)], fill=ICE_2)
    return img


def thunder_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon([(9, 2), (4, 9), (8, 9), (7, 14), (12, 7), (8, 7)], fill=THUNDER_1)
    d.polygon([(8, 4), (6, 8), (8, 8), (7, 11), (10, 7), (8, 7)], fill=THUNDER_2)
    return img


def score_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon(
        [(8, 3), (9, 7), (13, 8), (9, 9), (8, 13), (7, 9), (3, 8), (7, 7)],
        fill=THUNDER_1,
    )
    d.polygon([(8, 6), (10, 8), (8, 10), (6, 8)], fill=THUNDER_2)
    return img


def kills_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.rectangle((5, 4, 12, 10), fill=rgba("f0f4f8"))
    d.rectangle((6, 11, 11, 12), fill=rgba("f0f4f8"))
    d.rectangle((6, 7, 7, 8), fill=INK)
    d.rectangle((10, 7, 11, 8), fill=INK)
    d.rectangle((7, 11, 7, 12), fill=INK)
    d.rectangle((10, 11, 10, 12), fill=INK)
    return img


def distance_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon([(4, 4), (8, 8), (4, 12), (5, 12), (9, 8), (5, 4)], fill=ICE_1)
    d.polygon([(8, 4), (12, 8), (8, 12), (9, 12), (13, 8), (9, 4)], fill=ICE_2)
    return img


def best_icon():
    img = Image.new("RGBA", (16, 16), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.polygon(
        [(3, 5), (6, 8), (8, 4), (10, 8), (13, 5), (12, 10), (4, 10)], fill=THUNDER_1
    )
    d.rectangle((5, 11, 11, 12), fill=THUNDER_2)
    return img


def far_background():
    img = gradient(480, 270, SKY_TOP, SKY_MID, SKY_LOW)
    d = ImageDraw.Draw(img)
    d.ellipse((340, 18, 420, 98), fill=SUN)
    for x, y in [(44, 28), (104, 48), (152, 34), (211, 61), (284, 29), (438, 54)]:
        d.rectangle((x, y, x + 1, y + 1), fill=SUN)
    d.polygon(
        [
            (0, 176),
            (62, 122),
            (129, 167),
            (192, 118),
            (251, 170),
            (330, 108),
            (421, 165),
            (480, 132),
            (480, 270),
            (0, 270),
        ],
        fill=MOUNTAIN_FAR,
    )
    d.polygon(
        [
            (0, 197),
            (53, 152),
            (110, 188),
            (171, 142),
            (245, 201),
            (324, 136),
            (399, 191),
            (480, 151),
            (480, 270),
            (0, 270),
        ],
        fill=MOUNTAIN_MID,
    )
    return img


def mid_background():
    img = Image.new("RGBA", (480, 270), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    cloud_color = rgba("ffffff", 120)
    for x, y, w in [(34, 52, 84), (136, 68, 96), (274, 44, 72), (360, 76, 88)]:
        d.rounded_rectangle((x, y, x + w, y + 18), radius=9, fill=cloud_color)
    mesa = rgba("49636f", 180)
    d.polygon(
        [
            (0, 205),
            (46, 183),
            (96, 193),
            (132, 170),
            (175, 191),
            (230, 161),
            (302, 205),
            (356, 172),
            (420, 198),
            (480, 183),
            (480, 270),
            (0, 270),
        ],
        fill=mesa,
    )
    d.polygon(
        [
            (0, 225),
            (38, 210),
            (74, 218),
            (124, 196),
            (166, 216),
            (214, 190),
            (272, 228),
            (330, 196),
            (392, 220),
            (452, 204),
            (480, 214),
            (480, 270),
            (0, 270),
        ],
        fill=rgba("2b4054", 170),
    )
    return img


def near_background():
    img = Image.new("RGBA", (480, 270), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    fill = MOUNTAIN_NEAR
    for x in [26, 112, 208, 312, 406]:
        d.rectangle((x, 156, x + 10, 240), fill=fill)
        d.rectangle((x - 10, 240, x + 20, 248), fill=fill)
        d.rectangle((x + 3, 148, x + 7, 156), fill=fill)
    for x in [68, 170, 266, 364, 452]:
        d.rectangle((x, 176, x + 7, 248), fill=fill)
        d.polygon([(x - 12, 188), (x + 3, 164), (x + 18, 188)], fill=fill)
    d.rectangle((0, 248, 479, 270), fill=fill)
    return img


def make_sheet(frames):
    width = sum(frame.width for frame in frames)
    height = max(frame.height for frame in frames)
    sheet = Image.new("RGBA", (width, height), (0, 0, 0, 0))
    cursor = 0
    for frame in frames:
        sheet.paste(frame, (cursor, 0), frame)
        cursor += frame.width
    return sheet


def app_icon():
    img = Image.new("RGBA", (48, 48), SKY_MID)
    d = ImageDraw.Draw(img)
    d.rectangle((0, 36, 47, 47), fill=LEAF)
    d.rectangle((0, 42, 47, 47), fill=DIRT_MID)
    player = player_frame(0)
    img.paste(player, (0, 0), player)
    return img


def main():
    save("player.png", make_sheet([player_frame(i) for i in range(8)]))
    save(
        "enemies.png",
        make_sheet(
            [enemy_ground_frame(i) for i in range(3)]
            + [enemy_turret_frame(i) for i in range(3)]
            + [enemy_flyer_frame(i) for i in range(3)]
        ),
    )
    save(
        "projectiles.png",
        make_sheet(
            [
                projectile_fire(True),
                projectile_ice(True),
                projectile_thunder(True),
                projectile_fire(False),
                projectile_ice(False),
                projectile_thunder(False),
            ]
        ),
    )
    save(
        "tiles.png",
        make_sheet([tile_ground(), tile_platform(), tile_crate(), tile_spike()]),
    )
    save("fx.png", make_sheet([fire_fx(), ice_fx(), thunder_fx()]))
    save(
        "ui.png",
        make_sheet(
            [
                heart_icon(),
                fire_icon(),
                ice_icon(),
                thunder_icon(),
                score_icon(),
                kills_icon(),
                distance_icon(),
                best_icon(),
            ]
        ),
    )
    save("background_far.png", far_background())
    save("background_mid.png", mid_background())
    save("background_near.png", near_background())
    save("icon.png", app_icon())
    print(f"Generated assets in {OUT}")


if __name__ == "__main__":
    main()
