import torch
from PIL import Image
import open_clip


import open_clip
import torch




model, _, transform = open_clip.create_model_and_transforms('ViT-B-32', pretrained='laion2b_s34b_b79k')
model.cuda().eval()

new_captions = []
for path_name in path_names:
  im = Image.open(path_name).convert("RGB")
  im = transform(im).unsqueeze(0).cuda()
  with torch.no_grad(), torch.cuda.amp.autocast():
    generated = model.generate(im)
  caption = open_clip.decode(generated[0]).split("")[0].replace("", "")[:-2]
  new_captions.append(caption)