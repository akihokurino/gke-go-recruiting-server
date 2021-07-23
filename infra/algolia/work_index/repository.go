package work_index

import (
	"context"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/pkg/errors"

	"gke-go-sample/domain"
	pb "gke-go-sample/proto/go/pb"

	"gke-go-sample/infra/algolia"

	"gke-go-sample/adapter"
	"gke-go-sample/domain/work_domain"
)

func NewRepo(client *algolia.Client) adapter.WorkIndexRepo {
	return &repository{
		client: client,
	}
}

type repository struct {
	client *algolia.Client
}

func (r *repository) SearchWithPager(
	ctx context.Context,
	q string,
	businessCondition pb.Department_BusinessCondition,
	prefID domain.PrefID,
	mAreaID domain.MAreaID,
	sAreaID domain.SAreaID,
	railID domain.RailID,
	stationID domain.StationID,
	merit pb.Work_Merit,
	pager *domain.Pager,
	order pb.SearchWorkOrder) ([]domain.WorkID, error) {
	var indexes []index

	index := r.client.Index(newerIndexName)
	filter := r.genSearchFilter(businessCondition, prefID, mAreaID, sAreaID, railID, stationID, merit)
	switch order {
	case pb.SearchWorkOrder_HourPayment_Higher:
		index = r.client.Index(paymentHigherIndexName)
		filter = r.genSearchFilter(businessCondition, prefID, mAreaID, sAreaID, railID, stationID, merit)
	case pb.SearchWorkOrder_DayPayment_Higher:
		index = r.client.Index(paymentHigherIndexName)
		filter = r.genSearchFilter(businessCondition, prefID, mAreaID, sAreaID, railID, stationID, merit)
	}

	params := []interface{}{
		opt.HitsPerPage(pager.Limit()),
		opt.Page(pager.AlgoliaPage()),
		opt.Filters(filter),
	}

	searchRes, err := index.Search(q, params...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = searchRes.UnmarshalHits(&indexes); err != nil {
		return nil, errors.WithStack(err)
	}

	ids := make([]domain.WorkID, 0, len(indexes))
	for _, index := range indexes {
		ids = append(ids, domain.WorkID(index.ObjectID))
	}

	return ids, nil
}

func (r *repository) SearchCount(
	ctx context.Context,
	q string,
	businessCondition pb.Department_BusinessCondition,
	prefID domain.PrefID,
	mAreaID domain.MAreaID,
	sAreaID domain.SAreaID,
	railID domain.RailID,
	stationID domain.StationID,
	merit pb.Work_Merit,
	order pb.SearchWorkOrder) (uint64, error) {

	index := r.client.Index(newerIndexName)
	filter := r.genSearchFilter(businessCondition, prefID, mAreaID, sAreaID, railID, stationID, merit)
	switch order {
	case pb.SearchWorkOrder_HourPayment_Higher:
		index = r.client.Index(paymentHigherIndexName)
		filter = r.genSearchFilter(businessCondition, prefID, mAreaID, sAreaID, railID, stationID, merit)
	case pb.SearchWorkOrder_DayPayment_Higher:
		index = r.client.Index(paymentHigherIndexName)
		filter = r.genSearchFilter(businessCondition, prefID, mAreaID, sAreaID, railID, stationID, merit)
	}

	params := []interface{}{
		opt.Filters(filter),
	}

	searchRes, err := index.Search(q, params...)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return uint64(searchRes.NbHits), nil
}

func (r *repository) genSearchFilter(
	businessCondition pb.Department_BusinessCondition,
	prefID domain.PrefID,
	mAreaID domain.MAreaID,
	sAreaID domain.SAreaID,
	railID domain.RailID,
	stationID domain.StationID,
	merit pb.Work_Merit) string {
	str := fmt.Sprintf("status = %d", int32(pb.Work_Status_Active))

	if businessCondition != pb.Department_BusinessCondition_Unknown {
		str = fmt.Sprintf("%s AND businessCondition = %d", str, int32(businessCondition))
	}

	if prefID.String() != "" {
		str = fmt.Sprintf("%s AND prefId:%s", str, prefID.String())
	}

	if mAreaID.String() != "" {
		str = fmt.Sprintf("%s AND mAreaId:%s", str, mAreaID.String())
	}

	if sAreaID.String() != "" {
		str = fmt.Sprintf("%s AND sAreaId:%s", str, sAreaID.String())
	}

	if railID.String() != "" {
		str = fmt.Sprintf("%s AND railIds:%s", str, railID.String())
	}

	if stationID.String() != "" {
		str = fmt.Sprintf("%s AND stationIds:%s", str, stationID.String())
	}

	if merit != pb.Work_Merit_Unknown {
		str = fmt.Sprintf("%s AND merits = %d", str, int32(merit))
	}

	return str
}

func (r *repository) Save(ctx context.Context, work *work_domain.Work) error {
	if _, err := r.client.Index(newerIndexName).SaveObject(entityFrom(work)); err != nil {
		return errors.WithStack(err)
	}

	if _, err := r.client.Index(paymentHigherIndexName).SaveObject(entityFrom(work)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) SaveMulti(ctx context.Context, works []*work_domain.Work) error {
	indexes := make([]*index, 0, len(works))
	for _, w := range works {
		indexes = append(indexes, entityFrom(w))
	}

	if _, err := r.client.Index(newerIndexName).SaveObjects(indexes); err != nil {
		return errors.WithStack(err)
	}

	if _, err := r.client.Index(paymentHigherIndexName).SaveObjects(indexes); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, workID domain.WorkID) error {
	if _, err := r.client.Index(newerIndexName).DeleteObject(workID.String()); err != nil {
		return errors.WithStack(err)
	}

	if _, err := r.client.Index(paymentHigherIndexName).DeleteObject(workID.String()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) DeleteMulti(ctx context.Context, workIDs []domain.WorkID) error {
	ids := make([]string, 0, len(workIDs))
	for _, id := range workIDs {
		ids = append(ids, id.String())
	}

	if _, err := r.client.Index(newerIndexName).DeleteObjects(ids); err != nil {
		return errors.WithStack(err)
	}

	if _, err := r.client.Index(paymentHigherIndexName).DeleteObjects(ids); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
