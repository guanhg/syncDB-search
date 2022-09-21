package main

import (
	"context"
	"encoding/json"
	"fmt"
	schema "github/guanhg/syncDB-search/schema"

	"github.com/olivere/elastic/v7"
)

func testAggregation() {
	q := elastic.NewBoolQuery()
	q.Must(elastic.NewTermQuery("sm_id", 1716)).Must(elastic.NewTermQuery("medium_type", 0))
	q.Must(elastic.NewRangeQuery("medium_id").Gt(0))

	sumAgg := elastic.NewSumAggregation().Field("weight")
	disAgg := elastic.NewTermsAggregation().Field("medium_id").SubAggregation("weight", sumAgg).Size(50).OrderByAggregation("weight", false)

	search := schema.NewSearch(q, "sm_record_*")
	search.Size(0).Aggregation("track", disAgg)
	res, _ := search.Do(context.Background())
	aggResult, _ := res.Aggregations["track"].MarshalJSON()

	aggTrack := make(map[string]interface{})
	json.Unmarshal(aggResult, &aggTrack)

	fmt.Println(aggTrack)
}

func testAggregation2() {
	q := elastic.NewBoolQuery()
	q.Must(elastic.NewRangeQuery("id").Lte(53))

	sumAgg := elastic.NewSumAggregation().Field("weight")

	search := schema.NewSearch(q, "sm_record_2017")
	search.Size(10).Aggregation("s", sumAgg)
	res, _ := search.Do(context.Background())
	aggResult, _ := res.Aggregations["s"].MarshalJSON()

	aggTrack := make(map[string]interface{})
	json.Unmarshal(aggResult, &aggTrack)

	fmt.Println(aggTrack)
}

func testSearch() {
	q := elastic.NewRangeQuery("id").Gte(50).Lte(55)
	search := schema.NewSearch(q, "sm_record_2017")
	//search.Size(2).From(1)
	res := search.Result()
	hits := res["items"].([]interface{})
	hit := hits[0].(map[string]interface{})
	fmt.Println(res, hit)
}
