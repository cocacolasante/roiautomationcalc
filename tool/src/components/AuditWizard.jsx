import ProgressBar from './ui/ProgressBar';
import LiveROICounter from './ui/LiveROICounter';
import WelcomeStep from './steps/WelcomeStep';
import BusinessProfileStep from './steps/BusinessProfileStep';
import WorkflowStep from './steps/WorkflowStep';
import ContactStep from './steps/ContactStep';
import ProcessingStep from './steps/ProcessingStep';

function groupQuestionsByCategory(questions) {
  const groups = {};
  (questions || []).forEach((q) => {
    if (!groups[q.category]) groups[q.category] = [];
    groups[q.category].push(q);
  });
  return Object.entries(groups).map(([category, qs]) => ({ category, questions: qs }));
}

export default function AuditWizard({ config, step, answers, liveROI, onAnswer, onNext, onBack, onComplete }) {
  const categoryGroups = groupQuestionsByCategory(config.questions);

  const steps = [
    { component: WelcomeStep, key: 'welcome' },
    { component: BusinessProfileStep, key: 'profile' },
    ...categoryGroups.map((cat) => ({
      component: WorkflowStep,
      key: cat.category,
      questions: cat.questions,
      categoryLabel: cat.category.charAt(0).toUpperCase() + cat.category.slice(1),
    })),
    { component: ContactStep, key: 'contact' },
    { component: ProcessingStep, key: 'processing' },
  ];

  const totalSteps = steps.length;
  const progress = (step / Math.max(totalSteps - 1, 1)) * 100;
  const currentStep = steps[Math.min(step, steps.length - 1)];
  const StepComponent = currentStep.component;

  return (
    <div className="wizard-container">
      <div style={{ background: 'white', borderBottom: '1px solid #eee', padding: '1rem 1.5rem', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <div style={{ fontWeight: 700, fontSize: '1.1rem', color: 'var(--primary)' }}>{config.companyName}</div>
        <div style={{ fontSize: '0.875rem', color: '#888' }}>Step {Math.min(step + 1, totalSteps)} of {totalSteps}</div>
      </div>

      <ProgressBar progress={progress} />

      {liveROI && liveROI.recoverableCostPerYear > 500 && (
        <LiveROICounter value={liveROI.recoverableCostPerYear} />
      )}

      <div className="step-container">
        <StepComponent
          config={config}
          answers={answers}
          questions={currentStep.questions}
          categoryLabel={currentStep.categoryLabel}
          liveROI={liveROI}
          onAnswer={onAnswer}
          onNext={onNext}
          onBack={step > 0 ? onBack : null}
          onComplete={onComplete}
        />
      </div>
    </div>
  );
}
