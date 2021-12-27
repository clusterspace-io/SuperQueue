import json

with open('out.log') as f:
  lines = f.readlines()

  buckets = {
    "POSTRecord": {
      "3": 0,
      "6": 0,
      "9": 0,
      "15": 0,
      "30": 0,
      "inf": 0,
    },
    "GETRecord": {
      "3": 0,
      "6": 0,
      "9": 0,
      "15": 0,
      "30": 0,
      "inf": 0,
    },
    "Ack": {
      "3": 0,
      "6": 0,
      "9": 0,
      "15": 0,
      "30": 0,
      "inf": 0,
    },
  }

  for i in lines:
    if i.startswith("{"):
      # Found JSON (dirty)
      l = json.loads(i)
      if l["method"] == "POST" and l["uri"].startswith("/ack"):
        if l["latency"]<=3000000:
          buckets["Ack"]["3"] += 1
        if l["latency"]<=6000000:
          buckets["Ack"]["6"] += 1
        if l["latency"]<=9000000:
          buckets["Ack"]["9"] += 1
        if l["latency"]<=15000000:
          buckets["Ack"]["15"] += 1
        if l["latency"]<=30000000:
          buckets["Ack"]["30"] += 1
        if l["latency"]>30000000:
          buckets["Ack"]["inf"] += 1
      elif l["method"] == "POST" and l["uri"].startswith("/record"):
        if l["latency"]<=3000000:
          buckets["POSTRecord"]["3"] += 1
        if l["latency"]<=6000000:
          buckets["POSTRecord"]["6"] += 1
        if l["latency"]<=9000000:
          buckets["POSTRecord"]["9"] += 1
        if l["latency"]<=15000000:
          buckets["POSTRecord"]["15"] += 1
        if l["latency"]<=30000000:
          buckets["POSTRecord"]["30"] += 1
        if l["latency"]>30000000:
          buckets["POSTRecord"]["inf"] += 1
      elif l["method"] == "GET" and l["uri"].startswith("/record"):
        if l["latency"]<=3000000:
          buckets["GETRecord"]["3"] += 1
        if l["latency"]<=6000000:
          buckets["GETRecord"]["6"] += 1
        if l["latency"]<=9000000:
          buckets["GETRecord"]["9"] += 1
        if l["latency"]<=15000000:
          buckets["GETRecord"]["15"] += 1
        if l["latency"]<=30000000:
          buckets["GETRecord"]["30"] += 1
        if l["latency"]>30000000:
          buckets["GETRecord"]["inf"] += 1

  print(json.dumps(buckets))
