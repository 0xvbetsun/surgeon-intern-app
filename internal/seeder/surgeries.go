package seeder

import (
	"encoding/csv"
	"os"

	"github.com/gofrs/uuid"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"go.uber.org/zap"
)

type surgeryForSeeding struct {
	diagnose *models.SurgeryDiagnosis
	method   *models.SurgeryMethod
}

func getSurgeriesForSeeding() []*surgeryForSeeding {
	var surgeries = make([]*surgeryForSeeding, 0)
	records, err := readData("surgeries.csv")
	if err != nil {
		zap.S().Fatalf(err.Error())
	}
	for _, record := range records {
		if record[0] != "" && record[4] != "" {
			surgeries = append(surgeries, &surgeryForSeeding{
				diagnose: &models.SurgeryDiagnosis{
					ID:           uuid.Must(uuid.NewV4()).String(),
					Bodypart:     record[0],
					DiagnoseName: record[1],
					DiagnoseCode: record[2],
					ExtraCode:    record[3],
				},
				method: &models.SurgeryMethod{
					ID:           uuid.Must(uuid.NewV4()).String(),
					MethodName:   record[4],
					ApproachName: record[5],
					MethodCode:   record[6],
				},
			})
		} else if record[0] == "" && record[4] != "" {
			surgeries = append(surgeries, &surgeryForSeeding{
				diagnose: surgeries[len(surgeries)-1].diagnose,
				method: &models.SurgeryMethod{
					ID:           uuid.Must(uuid.NewV4()).String(),
					MethodName:   record[4],
					ApproachName: record[5],
					MethodCode:   record[6],
				},
			})
		}
	}
	zap.S().Info(surgeries)
	return surgeries
}
func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
