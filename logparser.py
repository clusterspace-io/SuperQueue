import json

with open('out.log') as f:
  lines = f.readlines()

  buckets = {
    "POSTRecord": {
      "le3": 0,
      "le6": 0,
      "le9": 0,
      "le15": 0,
      "le30": 0,
      "inf": 0,
    },
    "GETRecord": {
      "le3": 0,
      "le6": 0,
      "le9": 0,
      "le15": 0,
      "le30": 0,
      "inf": 0,
    },
    "Ack": {
      "le3": 0,
      "le6": 0,
      "le9": 0,
      "le15": 0,
      "le30": 0,
      "inf": 0,
    },
  }

  for i in lines:
    if i.startswith("{"):
      # Found JSON (dirty)
      l = json.loads(i)
      if l["method"] == "POST" and l["uri"].startswith("/ack"):
        if l["latency"]<=3000000:
          buckets["Ack"]["le3"] += 1
        if l["latency"]<=6000000:
          buckets["Ack"]["le6"] += 1
        if l["latency"]<=9000000:
          buckets["Ack"]["le9"] += 1
        if l["latency"]<=15000000:
          buckets["Ack"]["le15"] += 1
        if l["latency"]<=30000000:
          buckets["Ack"]["le30"] += 1
        if l["latency"]>30000000:
          buckets["Ack"]["inf"] += 1
      elif l["method"] == "POST" and l["uri"].startswith("/record"):
        if l["latency"]<=3000000:
          buckets["POSTRecord"]["le3"] += 1
        if l["latency"]<=6000000:
          buckets["POSTRecord"]["le6"] += 1
        if l["latency"]<=9000000:
          buckets["POSTRecord"]["le9"] += 1
        if l["latency"]<=15000000:
          buckets["POSTRecord"]["le15"] += 1
        if l["latency"]<=30000000:
          buckets["POSTRecord"]["le30"] += 1
        if l["latency"]>30000000:
          buckets["POSTRecord"]["inf"] += 1
      elif l["method"] == "GET" and l["uri"].startswith("/record"):
        if l["latency"]<=3000000:
          buckets["GETRecord"]["le3"] += 1
        if l["latency"]<=6000000:
          buckets["GETRecord"]["le6"] += 1
        if l["latency"]<=9000000:
          buckets["GETRecord"]["le9"] += 1
        if l["latency"]<=15000000:
          buckets["GETRecord"]["le15"] += 1
        if l["latency"]<=30000000:
          buckets["GETRecord"]["le30"] += 1
        if l["latency"]>30000000:
          buckets["GETRecord"]["inf"] += 1

  print(json.dumps(buckets))
