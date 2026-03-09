from dataclasses import dataclass, field

import xmlschema
from xmlschema.validators.groups import XsdGroup


@dataclass
class SchemaMetadata:
    children_order: dict[str, list[str]] = field(default_factory=dict)
    # dict where the keys are full paths to the node
    # and values are types of the node (e.g. TKwotaC)
    required_fields: dict[str, str] = field(default_factory=dict)
    # list of full paths to the array nodes
    array_nodes: set[str] = field(default_factory=set)


def schema_metadata(schema: xmlschema.XMLSchema) -> SchemaMetadata:
    meta = SchemaMetadata()

    def walk_element(xsd_element, path):
        xsd_type = xsd_element.type

        # Only complex types can have children
        if not xsd_type.is_complex():
            return

        model = xsd_type.content

        # Only groups (sequence / choice / all) contain elements
        if not isinstance(model, XsdGroup):
            return

        for child in model.iter_elements():
            # we'll decide what to do with required fields outside of this function
            # here let's just report it.
            if child.min_occurs > 0:
                meta.required_fields[f"{path}.{child.local_name}"] = (
                    child.type.local_name
                )

            meta.children_order.setdefault(path, []).append(child.local_name)

            # we need to know which elements are actually arrays - that is
            # necesary when we call Node.ToWriter function that actually
            # renders the content of XML
            if child.max_occurs is None or child.max_occurs > 1:
                meta.array_nodes.add(path + "." + child.local_name)

            walk_element(child, f"{path}.{child.local_name}")

    for element in schema.elements.values():
        walk_element(element, element.local_name)

    return meta


def load_schema_metadata(schema_path: str) -> SchemaMetadata:
    """Load and extract metadata from an XSD schema."""
    schema = xmlschema.XMLSchema(schema_path)
    return schema_metadata(schema)
