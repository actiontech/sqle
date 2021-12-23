import os
import sys

scspell_dir = []


def get_all_folders_containing_go_files(path):
    files = os.listdir(path)
    not_add_to_scspell_dir = True
    for file in files:
        if os.path.splitext(file)[-1] == ".go" and not_add_to_scspell_dir:
            scspell_dir.append(path)
            not_add_to_scspell_dir = False
        complete_path = path + "/" + file
        if os.path.isdir(complete_path):
            get_all_folders_containing_go_files(complete_path)


get_all_folders_containing_go_files(sys.argv[1])
d = os.system("scspell --use-builtin-base-dict --override-dictionary ./spelling_dict.txt --report-only " + "/*.go ".join(scspell_dir) + "/*.go")
print(d)
