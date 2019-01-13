import json
import os.path as osp


file_name = "sites.json"

def main():
    if not osp.isfile(file_name):
        raise FileNotFoundError()
    json_data: dict = json.loads(open(file_name, 'r', encoding="utf8").read())
    
    with open("sites.md", 'w', encoding="utf8") as md:
        md.write("# {n} sites are supported!\n".format(n=len(json_data)))
        for name in json_data.keys():
            md.write(" - [{sns_name}]({url})\n".format(sns_name=name, url=str(
                json_data[name]).replace('?', str()).replace("?.", str())))

if __name__ == "__main__":
    exit(
        main()
    )