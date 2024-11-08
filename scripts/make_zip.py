import os, hashlib, zipfile, sys
root = os.path.dirname(os.path.dirname(__file__))
zip_path = os.path.join(root, "A-secrets-infra-starter.zip")
with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as z:
    for folder, _, files in os.walk(root):
        for f in files:
            if f.endswith(".zip"): continue
            fp = os.path.join(folder, f)
            arc = os.path.relpath(fp, root)
            z.write(fp, arc)
h = hashlib.sha256()
with open(zip_path, "rb") as fh:
    h.update(fh.read())
print(zip_path)
print("sha256:", h.hexdigest())
