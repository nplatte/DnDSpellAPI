import unittest
from parse import Spell


class TestSpell(unittest.TestCase):

    def setUp(self):
        self.spell = Spell()
        return super().setUp()
    
    def tearDown(self):
        return super().tearDown()
    
    def test_set_spell_name(self):
        test_string = "Firebolt"
        self.assertEqual(self.spell.name, "")
        self.spell.set_name(test_string)
        self.assertEqual(self.spell.name, test_string)

    def test_set_spell_range(self):
        test_string = "Range: 120 feet"
        self.assertEqual(self.spell.range, "")
        self.spell.set_range(test_string)
        self.assertEqual(self.spell.range, "120 feet")

    def test_set_spell_school(self):
        test_string = "Conjuration"
        self.assertEqual(self.spell.school, "")
        self.spell.set_school(test_string)
        self.assertEqual(self.spell.school, test_string)

    def test_set_spell_level(self):
        test_string = "cantrip"
        self.assertEqual(self.spell.level, "")
        self.spell.set_level(test_string)
        self.assertEqual(self.spell.level, test_string)

    def test_set_cast_time(self):
        test_string = "Casting Time: 1 action"
        self.assertEqual(self.spell.cast_time, "")
        self.spell.set_cast_time(test_string)
        self.assertEqual(self.spell.cast_time, "1 action")

    def test_set_cast_time_reaction(self):
        test_string = "Casting Time: 1 reaction, which you take when you are damaged by a creature within 60 feet of you that you can see"
        self.assertEqual(self.spell.cast_time, "")
        self.spell.set_cast_time(test_string)
        self.assertEqual(self.spell.cast_time, "1 reaction")

    def test_set_components_V(self):
        test_string = "Components: V"
        self.assertEqual(self.spell.components, [])
        self.spell.set_components(test_string)
        self.assertIn("V", self.spell.components)
        self.assertEqual(len(self.spell.components), 1)

    def test_set_components_V_S_M(self):
        test_string = "Components: V, S, M (mistletoe, a shamrock leaf, and a club or quarterstaff)"
        self.assertEqual(self.spell.components, [])
        self.spell.set_components(test_string)
        self.assertIn("V", self.spell.components)
        self.assertIn("S", self.spell.components)
        self.assertIn("M", self.spell.components)
        self.assertIn("mistletoe, a shamrock leaf, and a club or quarterstaff", self.spell.components)
        self.assertEqual(len(self.spell.components), 4)

    def test_set_duration(self):
        test_string = "Duration: Instantaneous"
        self.assertEqual(self.spell.duration, "")
        self.spell.set_duration(test_string)
        self.assertEqual(self.spell.duration, "Instantaneous")

    def test_set_class_list_one(self):
        test_string = "Spell Lists. Druid"
        self.assertEqual(len(self.spell.class_list), 0)
        self.spell.set_class_list(test_string)
        self.assertIn("Druid", self.spell.class_list)
        self.assertEqual(len(self.spell.class_list), 1)

    def test_set_class_list_five(self):
        test_string = "Spell Lists. Artificer, Bard, Cleric, Druid, Sorcerer, Wizard"
        self.assertEqual(len(self.spell.class_list), 0)
        self.spell.set_class_list(test_string)
        self.assertIn("Artificer", self.spell.class_list)
        self.assertIn("Bard", self.spell.class_list)
        self.assertIn("Cleric", self.spell.class_list)
        self.assertIn("Druid", self.spell.class_list)
        self.assertIn("Sorcerer", self.spell.class_list)
        self.assertIn("Wizard", self.spell.class_list)
        self.assertEqual(len(self.spell.class_list), 6)

    def test_set_higher_levels_true(self):
        test_string = "At Higher Levels. This spell’s damage increases by 1d6 when you reach 5th level (2d6), 11th level (3d6), and 17th level (4d6)."
        self.assertEqual(self.spell.higher_levels, "")
        self.spell.set_higher_levels(test_string)
        self.assertEqual(self.spell.higher_levels, "This spell’s damage increases by 1d6 when you reach 5th level (2d6), 11th level (3d6), and 17th level (4d6).")

    def test_set_higher_levels_false(self):
        test_string = "You instantly light or snuff out a candle, a torch, or a small campfire."
        self.assertEqual(self.spell.higher_levels, "")
        self.spell.set_higher_levels(test_string)
        self.assertEqual(self.spell.higher_levels, "")

    def test_set_description_higher(self):
        test_list = ['You hurl a bubble of acid. ', 'Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take 1d6 acid damage.', 'At Higher Levels. This spell’s damage increases by 1d6 when you reach 5th level (2d6), 11th level (3d6), and 17th level (4d6).']
        self.assertEqual(self.spell.description, "")
        self.spell.set_description(test_list)
        correct = 'You hurl a bubble of acid. Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take 1d6 acid damage.'
        self.assertEqual(self.spell.description, correct)

    def test_set_description_no_higher(self):
        test_list = ['You hurl a bubble of acid. ', 'Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take 1d6 acid damage.']
        self.assertEqual(self.spell.description, "")
        self.spell.set_description(test_list)
        correct = 'You hurl a bubble of acid. Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take 1d6 acid damage.'
        self.assertEqual(self.spell.description, correct)

    def test_set_ritual_true(self):
        test_string = "1st-level divination (ritual)"
        self.assertEqual(self.spell.ritual, "False")
        self.spell.set_ritual(test_string)
        self.assertEqual(self.spell.ritual, "True")

    def test_set_ritual_false(self):
        test_string = "1st-level divination"
        self.assertEqual(self.spell.ritual, "False")
        self.spell.set_ritual(test_string)
        self.assertEqual(self.spell.ritual, "False")

    def test_set_concentration_true(self):
        test_string = "Duration: Concentration, up to 1 minute"
        self.assertEqual(self.spell.concentration, "False")
        self.spell.set_concentration(test_string)
        self.assertEqual(self.spell.concentration, "True")

    def test_set_concentration_false(self):
        test_string = "Duration: up to 1 minute"
        self.assertEqual(self.spell.concentration, "False")
        self.spell.set_concentration(test_string)
        self.assertEqual(self.spell.concentration, "False")


if __name__ == "__main__":
    unittest.main()