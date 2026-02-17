from typing import Dict
import json

site_dict: Dict[str, Dict] = json.loads(open("data.json", "r", encoding="utf-8").read())

with open("sites.md", "w", encoding="utf-8") as file:
    file.write(f"# {len(site_dict)} sites are supported\n\n")
    for site in sorted(site_dict.keys()):
        if site != "$schema":
            file.write(f"- [{site}]({site_dict[site]['urlMain']})\n")

    file.write("\n")
    file.write(
        "## Removed sites\n\n"
        "Please refer to the [Sherlock database](https://github.com/sherlock-project/sherlock/blob/master/removed_sites.md)"
        "\n"
    )
