import json
import re
from os.path import isfile

file_name = "sites.json"


def main():
    if not isfile(file_name):
        raise FileNotFoundError("JSON data not found!")
    json_data: dict = json.loads(open(file_name, 'r', encoding="utf8").read())

    with open("sites.md", 'w', encoding="utf8") as md:
        md.write("# {n} sites are supported!\n".format(n=len(json_data)))
        for name in json_data.keys():
            s = re.search(
                r'(?P<http>https?://)(?P<www>www\.)?(\?\.)?(?P<main>(\.?\w+\.?)+)',
                json_data[name]
            )
            http = s.group("http")
            www = s.group("www")
            url_main = s.group("main")
            md.write(
                " - [{sns_name}]({url})\n".format(
                    sns_name=name,
                    url=(http if http else '') + (www if www else '') + url_main
                )
            )


if __name__ == "__main__":
    main()
