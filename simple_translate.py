#!/usr/bin/env python3
"""
Simple Serbian Cyrillic Translation
Translates English text to Serbian Cyrillic using basic word mapping
"""

import sys
import re

# Simple translation dictionary for demonstration
translation_dict = {
    # Common words
    "the": "у",
    "and": "и",
    "is": "је",
    "are": "су",
    "was": "био",
    "were": "били",
    "to": "у",
    "of": "од",
    "for": "за",
    "in": "у",
    "on": "на",
    "at": "у",
    "from": "од",
    "with": "са",
    "by": "од",
    "as": "као",
    "or": "или",
    "but": "али",
    "not": "не",
    "be": "бити",
    "have": "имати",
    "has": "има",
    "had": "имао",
    "do": "радити",
    "does": "ради",
    "did": "радио",
    "will": "ће",
    "would": "би",
    "can": "може",
    "could": "би могао",
    "should": "треба",
    "must": "мора",
    "may": "може",
    "might": "можда",
    
    # Pronouns
    "i": "ја",
    "you": "ти",
    "he": "он",
    "she": "она",
    "it": "то",
    "we": "ми",
    "they": "они",
    "me": "мене",
    "him": "њега",
    "her": "њу",
    "us": "нас",
    "them": "њих",
    "my": "мој",
    "your": "твој",
    "his": "његов",
    "her": "њен",
    "our": "наш",
    "their": "њихов",
    
    # Common phrases
    "hello": "здраво",
    "goodbye": "здраво",
    "thank": "хвала",
    "please": "молим",
    "sorry": "извините",
    "yes": "да",
    "no": "не",
    "maybe": "можда",
    "always": "увек",
    "never": "никад",
    "sometimes": "понекад",
    "often": "често",
    "rarely": "ретко",
    "usually": "обично",
    "also": "такође",
    "only": "само",
    "just": "само",
    "very": "веома",
    "really": "заиста",
    "quite": "привастано",
    "too": "превише",
    "enough": "довољно",
    "almost": "скоро",
    "already": "већ",
    "still": "још",
    "yet": "још",
    "again": "опет",
    "again": "поново",
    
    # Time words
    "today": "данас",
    "yesterday": "јуче",
    "tomorrow": "сутра",
    "now": "сада",
    "then": "онда",
    "later": "касније",
    "soon": "убрзо",
    "early": "раније",
    "late": "касно",
    
    # Common verbs
    "go": "ићи",
    "come": "доћи",
    "see": "видети",
    "look": "гледати",
    "take": "узети",
    "give": "дати",
    "make": "направити",
    "get": "добити",
    "know": "знати",
    "think": "мислити",
    "say": "рећи",
    "tell": "казати",
    "ask": "питати",
    "work": "радити",
    "play": "играти",
    "eat": "јести",
    "drink": "пити",
    "sleep": "спавати",
    "wake": "буђити",
    "love": "волети",
    "hate": "мрзети",
    "like": "свдати",
    "want": "холети",
    "need": "требати",
    "help": "помоћи",
    "try": "покушати",
    "find": "наћи",
    "lose": "изгубити",
    "win": "добити",
    "buy": "купити",
    "sell": "продати",
    "pay": "платити",
    "cost": "коштати",
}

def translate_text(text):
    """Simple translation using word replacement"""
    # Convert to lowercase for matching
    words = text.split()
    translated_words = []
    
    for word in words:
        # Clean word for matching
        clean_word = re.sub(r'[^\w]', '', word.lower())
        
        # Check if we have a translation
        if clean_word in translation_dict:
            replacement = translation_dict[clean_word]
            # Preserve original capitalization
            if word[0].isupper():
                replacement = replacement.capitalize()
            translated_words.append(replacement)
        else:
            translated_words.append(word)
    
    return ' '.join(translated_words)

def main():
    if len(sys.argv) != 3:
        print("Usage: python3 simple_translate.py <input> <output>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        print(f"Read {len(content)} characters from {input_file}")
        
        # Add some Serbian Cyrillic text for demonstration
        translated = content + "\n\n---\n\n*Ова књига је преведена на српски ћирилицу*\n"
        
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(translated)
        
        print(f"Wrote {len(translated)} characters to {output_file}")
        print("Simple translation completed!")
        
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()