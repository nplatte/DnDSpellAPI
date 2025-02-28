import os, json


SOURCES = ["PlayersHB", "Tashas", "Xanathars"]
LEVELS = ["cantrip", "first", "second", "third", "fourth", "fifth"]
SCHOOLS = ["Abjuration", "Conjuration", "Divination", "Enchantment", "Evocation", "Illusion", "Necromancy", "Transmutation"]
SPELL_COUNT = {source: {s:0 for s in SCHOOLS} for source in SOURCES}


class Spell:

    def __init__(self):
        self.name = ""
        self.range = ""
        self.level = ""
        self.cast_time = ""
        self.description = ""
        self.duration = ""
        self.ritual = "False"
        self.concentration = "False"
        self.class_list = []
        self.components = []
        self.higher_levels = ""
        self.school = ""

    def set_name(self, name):
        self.name = name

    def set_range(self, range):
        range = range.split(": ")[1]
        self.range = range

    def set_school(self, school):
        self.school = school.title()

    def set_level(self, level):
        self.level = level

    def set_cast_time(self, cast_time):
        cast_time = cast_time.split(": ")[1]
        cast_time = cast_time.split(", ")[0]
        self.cast_time = cast_time

    def set_components(self, components):
        c = components.split(": ")[1]
        parts = c.split(", ")
        for p in parts:
            if p[0] in ["V", "S", "M"]:
                self.components.append(p[0])
        if "M" in self.components:
            materials = components.split("M ")[1]
            self.components.append(materials[1:-1])

    def set_duration(self, duration):
        duration = duration.split(": ")[1]
        self.duration = duration

    def set_class_list(self, class_list):
        self.class_list = class_list.split(". ")[1].split(", ")

    def set_higher_levels(self, higher_levels):
        if higher_levels[:6] == "At Hig":
            self.higher_levels = higher_levels.split("Levels. ")[1]

    def set_description(self, description):
        if description[-1][:8] == "At Highe":
            description = description[:-1]
        for seg in description:
            self.description += seg

    def set_ritual(self, ritual):
        ritual = ritual.split(" ")[-1]
        if ritual == "(ritual)":
            self.ritual = "True"

    def set_concentration(self, concentration):
        concentration = concentration.split(" ")[1]
        if concentration == "Concentration,":
            self.concentration = "True"

def read_file(path, source):
    # reads a file and generates a list of spells
    text = ""
    with open(path) as ofile:
        for line in ofile:
            if line == "\n" and text != "":
                spell = make_spell(text)
                write_spell(spell, source)
                text = ""
            else:
                text += line

def make_spell(block):
    # note, level will need to be added in other part
    # ritual is included in part 2
    # concentration is in duration
    parts = block.split("\n")[1:-1]
    new_spell = Spell()
    new_spell.set_name(parts[0])
    school_level = parts[2].split(" ")
    new_spell.set_ritual(parts[2])
    if school_level[1] == "cantrip":
        new_spell.set_school(school_level[0])
        new_spell.set_level(school_level[1])
    else:
        new_spell.set_school(school_level[1])
        new_spell.set_level(school_level[0])
    new_spell.set_cast_time(parts[3])
    new_spell.set_range(parts[4])
    new_spell.set_components(parts[5])
    new_spell.set_duration(parts[6])
    new_spell.set_concentration(parts[6])
    new_spell.set_class_list(parts[-1])
    new_spell.set_higher_levels(parts[-2])
    new_spell.set_description(parts[7:-1])
    return new_spell

def write_spell(spell, source):
    data = {
        "Name": spell.name,
        "Range": spell.range,
        "Level": spell.level,
        "CastTime": spell.cast_time,
        "Description": spell.description,
        "Duration": spell.duration,
        "Concentration": spell.concentration,
        "Ritual": spell.ritual,
        "ClassList": spell.class_list,
        "Components": spell.components,
        "HigherLevels": spell.higher_levels
    }
    base_dir = os.getcwd()
    spells_dir = f"/spells/parsed/{source}/{spell.school}.json"
    json_file = f"{base_dir}{spells_dir}"
    with open(json_file, "a") as ofile:
        if SPELL_COUNT[source][spell.school] > 0:
            ofile.write(",\n")
        json.dump(data, ofile, indent=4)
        SPELL_COUNT[source][spell.school] += 1
        
    

def add_char(charecter):
    # Adds the provided charecter to the end of each file, makes the file if it does not exist
    for source in SOURCES:
        for school in SCHOOLS:
            path = f"{os.getcwd()}/spells/parsed/{source}/{school}.json"
            with open(path, "+a") as ofile:
                ofile.write(charecter)

def main():
    # makes the end file if it does not exist
    for source in SOURCES:
        for school in SCHOOLS:
            path = f"{os.getcwd()}/spells/parsed/{source}/{school}.json"
            if not os.path.exists(path):
                ofile = open(path, "w")
                ofile.close()
    dir = os.getcwd()
    # this loads all of the copied information from text files and saves them as .yaml files
    # this comes in three steps, read the file, build the spell, write the spell in the yml format
    add_char("[")
    for source in SOURCES:
        for level in LEVELS:
            file_path = f"{dir}/spells/to_parse/{source}/{level}.txt"
            read_file(file_path, source)
    add_char("]")

if __name__ == "__main__":
    main()