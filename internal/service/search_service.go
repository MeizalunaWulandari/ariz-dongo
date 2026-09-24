package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/client/fieldsa"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/logger"
)

type SearchService struct {
	fieldsaClient *fieldsa.Client
	detailService *DetailService
	concurrency   int
}

func NewSearchService(
	fieldsaClient *fieldsa.Client,
	detailService *DetailService,
	concurrency int,
) *SearchService {

	if concurrency < 1 {
		concurrency = 5
	}

	return &SearchService{
		fieldsaClient: fieldsaClient,
		detailService: detailService,
		concurrency:   concurrency,
	}
}

type searchResult struct {
	workOrderNumber string
	found           bool
}

func (s *SearchService) Search(
	ctx context.Context,
	circuitID string,
) (*WorkOrderDetailResponse, error) {

	circuitID = strings.TrimSpace(circuitID)

	if circuitID == "" {
		return nil, fmt.Errorf(
			"circuitId wajib diisi",
		)
	}

	year, month, err := parseCircuitID(circuitID)

	if err != nil {
		return nil, err
	}

	logger.Info(
		"Search circuitId %s: tahun=%d bulan=%02d",
		circuitID,
		year,
		month,
	)

	startDate := time.Date(
		year,
		month,
		1,
		0,
		0,
		0,
		0,
		time.Local,
	)

	daysInMonth := time.Date(
		year,
		month+1,
		0,
		0,
		0,
		0,
		0,
		time.Local,
	).Day()

	/*
		==================================================
		CANCELLATION
		==================================================

		ctxSearch akan dibatalkan segera setelah circuitId
		ditemukan.
	*/

	searchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	/*
		==================================================
		SHARED RESULT
		==================================================
	*/

	var (
		mu     sync.Mutex
		result searchResult
	)

	/*
		==================================================
		WORKER POOL
		==================================================

		Tidak membuat 28/30/31 goroutine sekaligus.

		Misal concurrency = 5:

		Worker 1 → tanggal
		Worker 2 → tanggal
		Worker 3 → tanggal
		Worker 4 → tanggal
		Worker 5 → tanggal

		Setelah ada yang menemukan data:
		cancel()
	*/

	jobs := make(chan time.Time)

	var wg sync.WaitGroup

	worker := func() {

		defer wg.Done()

		for date := range jobs {

			if searchCtx.Err() != nil {
				return
			}

			dateString := date.Format(
				"2006-01-02",
			)

			data, err :=
				s.fieldsaClient.GetWorkOrderDispatch(
					searchCtx,
					dateString,
					"MSwyLDMsNSw2LDc=",
					"KALIMANTAN",
				)

			if err != nil {

				if searchCtx.Err() != nil {
					return
				}

				logger.Warn(
					"Search %s tanggal %s gagal: %v",
					circuitID,
					dateString,
					err,
				)

				continue
			}

			workOrderNumber, found :=
				findCircuitID(
					data,
					circuitID,
				)

			if !found {
				continue
			}

			/*
				==================================================
				FOUND
				==================================================
			*/

			mu.Lock()

			if !result.found {

				result = searchResult{
					workOrderNumber: workOrderNumber,
					found:           true,
				}

				logger.Info(
					"Search %s FOUND: %s pada %s",
					circuitID,
					workOrderNumber,
					dateString,
				)

				cancel()
			}

			mu.Unlock()

			return
		}
	}

	for i := 0; i < s.concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	/*
		==================================================
		PRODUCE DATES
		==================================================
	*/

produce:
	for day := 0; day < daysInMonth; day++ {

		if searchCtx.Err() != nil {
			break produce
		}

		date := startDate.AddDate(
			0,
			0,
			day,
		)

		select {

		case jobs <- date:

		case <-searchCtx.Done():
			break produce
		}
	}

	close(jobs)

	wg.Wait()

	/*
		==================================================
		RESULT
		==================================================
	*/

	mu.Lock()
	foundWO := result.workOrderNumber
	found := result.found
	mu.Unlock()

	if !found {
		return nil, fmt.Errorf(
			"circuitId %s tidak ditemukan pada bulan %04d-%02d",
			circuitID,
			year,
			month,
		)
	}

	/*
		==================================================
		GET DETAIL
		==================================================

		Setelah WO ditemukan, kita menggunakan service
		detail yang sama dengan endpoint detail biasa.

		Artinya response search dan detail IDENTIK.
	*/

	return s.detailService.GetDetail(
		ctx,
		foundWO,
	)
}

func parseCircuitID(
	circuitID string,
) (int, time.Month, error) {

	/*
		Format:

		CRT2510014054
		    │││
		    ││└── tanggal
		    │└─── bulan
		    └──── tahun

		Contoh:

		CRT2510014054
		    25 = 2025
		    10 = Oktober
		    01 = tanggal

		CRT2609014054
		    26 = 2026
		    09 = September
		    01 = tanggal
	*/

	if len(circuitID) < 9 {
		return 0, 0, fmt.Errorf(
			"format circuitId tidak valid",
		)
	}

	if !strings.HasPrefix(
		strings.ToUpper(circuitID),
		"CRT",
	) {
		return 0, 0, fmt.Errorf(
			"circuitId harus diawali CRT",
		)
	}

	yearPart := circuitID[3:5]
	monthPart := circuitID[5:7]

	yearShort, err :=
		strconv.Atoi(yearPart)

	if err != nil {
		return 0, 0, fmt.Errorf(
			"tahun circuitId tidak valid",
		)
	}

	monthNumber, err :=
		strconv.Atoi(monthPart)

	if err != nil {
		return 0, 0, fmt.Errorf(
			"bulan circuitId tidak valid",
		)
	}

	if monthNumber < 1 || monthNumber > 12 {
		return 0, 0, fmt.Errorf(
			"bulan circuitId tidak valid: %d",
			monthNumber,
		)
	}

	return 2000 + yearShort,
		time.Month(monthNumber),
		nil
}

func findCircuitID(
	data []byte,
	circuitID string,
) (string, bool) {

	/*
		Dispatch response dari Fieldsa kita decode secara
		dinamis karena struktur response list dapat berubah.

		Kita hanya membutuhkan:

		circuitId
		workOrderNumber
	*/

	var raw any

	if err := json.Unmarshal(data, &raw); err != nil {

		logger.Warn(
			"Decode dispatch response gagal: %v",
			err,
		)

		return "", false
	}

	return findCircuitInValue(
		raw,
		circuitID,
	)
}

func findCircuitInValue(
	value any,
	circuitID string,
) (string, bool) {

	switch item := value.(type) {

	case []any:

		for _, child := range item {

			workOrderNumber, found :=
				findCircuitInValue(
					child,
					circuitID,
				)

			if found {
				return workOrderNumber, true
			}
		}

	case map[string]any:

		var (
			currentCircuit string
			workOrder      string
		)

		for key, value := range item {

			switch strings.ToLower(key) {

			case "circuitid":
				if text, ok := value.(string); ok {
					currentCircuit = text
				}

			case "workordernumber":
				if text, ok := value.(string); ok {
					workOrder = text
				}
			}
		}

		if strings.EqualFold(
			currentCircuit,
			circuitID,
		) && workOrder != "" {

			return workOrder, true
		}

		for _, child := range item {

			workOrderNumber, found :=
				findCircuitInValue(
					child,
					circuitID,
				)

			if found {
				return workOrderNumber, true
			}
		}
	}

	return "", false
}
