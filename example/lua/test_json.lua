function processor(key,count)
	data = {}
	data["hello"]="world"
	data["a"] = {}
	data["a"]["b"] = "b"
	data["a"]["c"] = {1,2,3,4,5,6}
	res = jsonMarshal(data)
	res = jsonUnMarshal(res)
	for k,v in ipairs(res["a"]["c"]) do
        print(k,v)
	end
	return true
end