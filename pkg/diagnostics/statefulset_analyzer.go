package diagnostics

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// StatefulSetAnalyzer provides analysis for StatefulSet resources
type StatefulSetAnalyzer struct {
	client    kubernetes.Interface
	namespace string
}

// NewStatefulSetAnalyzer creates a new StatefulSetAnalyzer
func NewStatefulSetAnalyzer(client kubernetes.Interface, namespace string) *StatefulSetAnalyzer {
	return &StatefulSetAnalyzer{
		client:    client,
		namespace: namespace,
	}
}

// StatefulSetReport contains the analysis report
type StatefulSetReport struct {
	Name                 string
	Namespace            string
	DesiredReplicas      int32
	CurrentReplicas      int32
	ReadyReplicas        int32
	UpdatedReplicas      int32
	Conditions           []appsv1.StatefulSetCondition
	PodTemplate          corev1.PodTemplateSpec
	VolumeClaimTemplates []corev1.PersistentVolumeClaim
	Events               []corev1.Event
	Analysis             StatefulSetAnalysis
}

// StatefulSetAnalysis contains diagnostic results
type StatefulSetAnalysis struct {
	Status          string
	Issues          []string
	Recommendations []string
	UpdateStrategy  string
}

// Analyze performs the analysis of a StatefulSet
func (s *StatefulSetAnalyzer) Analyze(statefulSetName string) (*StatefulSetReport, error) {
	statefulSet, err := s.client.AppsV1().StatefulSets(s.namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get statefulset %s: %v", statefulSetName, err)
	}

	events, err := s.client.CoreV1().Events(s.namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + statefulSetName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get events for statefulset %s: %v", statefulSetName, err)
	}

	report := &StatefulSetReport{
		Name:                 statefulSet.Name,
		Namespace:            statefulSet.Namespace,
		DesiredReplicas:      *statefulSet.Spec.Replicas,
		CurrentReplicas:      statefulSet.Status.Replicas,
		ReadyReplicas:        statefulSet.Status.ReadyReplicas,
		UpdatedReplicas:      statefulSet.Status.UpdatedReplicas,
		Conditions:           statefulSet.Status.Conditions,
		PodTemplate:          statefulSet.Spec.Template,
		VolumeClaimTemplates: statefulSet.Spec.VolumeClaimTemplates,
		Events:               events.Items,
	}

	s.analyzeConditions(report)
	s.analyzeUpdateStrategy(report, statefulSet)
	s.analyzeReplicaStatus(report)

	return report, nil
}

func (s *StatefulSetAnalyzer) analyzeConditions(report *StatefulSetReport) {
	for _, condition := range report.Conditions {
		// Check for any condition that indicates a problem
		if condition.Status == corev1.ConditionFalse {
			report.Analysis.Issues = append(report.Analysis.Issues,
				fmt.Sprintf("Condition %s is False: %s", condition.Type, condition.Message))
		}
	}
}

func (s *StatefulSetAnalyzer) analyzeReplicaStatus(report *StatefulSetReport) {
	if report.ReadyReplicas != report.DesiredReplicas {
		report.Analysis.Issues = append(report.Analysis.Issues,
			fmt.Sprintf("Ready replicas (%d) does not match desired replicas (%d)",
				report.ReadyReplicas, report.DesiredReplicas))
	}

	if report.CurrentReplicas != report.DesiredReplicas {
		report.Analysis.Issues = append(report.Analysis.Issues,
			fmt.Sprintf("Current replicas (%d) does not match desired replicas (%d)",
				report.CurrentReplicas, report.DesiredReplicas))
	}

	if len(report.Analysis.Issues) == 0 {
		report.Analysis.Status = "Healthy"
	} else {
		report.Analysis.Status = "Unhealthy"
	}
}

func (s *StatefulSetAnalyzer) analyzeUpdateStrategy(report *StatefulSetReport, statefulSet *appsv1.StatefulSet) {
	if statefulSet.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType {
		report.Analysis.UpdateStrategy = "RollingUpdate"
		if statefulSet.Spec.UpdateStrategy.RollingUpdate != nil &&
			statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition != nil {
			report.Analysis.Recommendations = append(report.Analysis.Recommendations,
				fmt.Sprintf("Partitioned update configured at partition %d",
					*statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition))
		}
	} else {
		report.Analysis.UpdateStrategy = "OnDelete"
		report.Analysis.Recommendations = append(report.Analysis.Recommendations,
			"Consider using RollingUpdate strategy for automated pod updates")
	}
}
