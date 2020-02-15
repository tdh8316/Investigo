from typing import Dict, AnyStr
import json

site_dict: Dict[AnyStr, AnyStr] = json.loads(
    open("data.json", "r", encoding="utf-8").read()
)

with open("sites.md", "w", encoding="utf-8") as file:
    file.write(f"# {len(site_dict)} sites are supprted!\n")
    for site in sorted(site_dict.keys()):
        file.write(f" - [{site}]({site_dict[site]['urlMain']})\n")
