import { Check } from 'lucide-react';

export interface StepperStep {
  key: string;
  label: string;
}

export interface StepperProps {
  steps: readonly StepperStep[];
  activeStep: number;
  onStepClick?: (index: number) => void;
}

export function Stepper({ steps, activeStep, onStepClick }: StepperProps) {
  return (
    <div className="w-[300px] mx-auto py-4">
      <div className="flex items-center justify-between relative">
        {/* Background Line */}
        <div className="absolute top-1/2 left-0 w-full h-0.5 bg-[var(--bg-muted)] -translate-y-1/2 z-0" />
        
        {/* Active Progress Line */}
        <div 
          className="absolute top-1/2 left-0 h-0.5 bg-[var(--primary)] -translate-y-1/2 z-0 transition-all duration-300 ease-in-out" 
          style={{ width: `${(activeStep / (steps.length - 1)) * 100}%` }}
        />

        {steps.map((step, index) => {
          const isCompleted = index < activeStep;
          const isActive = index === activeStep;
          const isStepDisabled = index > activeStep;

          return (
            <div key={step.key} className="relative z-10 flex flex-col items-center group">
              <button
                type="button"
                disabled={isStepDisabled}
                onClick={() => onStepClick?.(index)}
                className={`
                  w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold transition-all duration-300
                  ${isCompleted ? 'bg-[var(--primary)] text-white' : isActive ? 'bg-[var(--bg-elevated)] border-2 border-[var(--primary)] text-[var(--primary)] scale-110 shadow-md' : 'bg-[var(--bg-muted)] text-[var(--text-muted)] border-2 border-transparent'}
                  ${!isStepDisabled ? 'hover:shadow-lg cursor-pointer' : 'cursor-default'}
                `}
              >
                {isCompleted ? <Check size={16} /> : index + 1}
              </button>
              
              <span className={`
                absolute top-10 whitespace-nowrap text-[10px] sm:text-xs font-semibold tracking-tight transition-colors duration-300
                ${isActive ? 'text-[var(--primary)]' : isCompleted ? 'text-[var(--text)]' : 'text-[var(--text-muted)]'}
              `}>
                {step.label}
              </span>
            </div>
          );
        })}
      </div>
      {/* Spacer for labels */}
      <div className="h-8" />
    </div>
  );
}
