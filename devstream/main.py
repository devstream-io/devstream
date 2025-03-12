from typing import Optional

from .workflow import Workflow


def call(workflow_name: str, user_input: Optional[str] = None):
    user_input = user_input or ""

    wf_name, _ = Workflow.parse_trigger(workflow_name)
    workflow = Workflow.load(wf_name) if wf_name else None

    if not workflow:
        raise ValueError(f"Workflow <{workflow_name}> not found.")

    if workflow.should_show_help(user_input):
        doc = workflow.get_help_doc(user_input)
        return doc

    workflow.setup(user_input, None, None, None)

    return workflow.run()
