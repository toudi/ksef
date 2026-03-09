from pathlib import Path
from subprocess import call

from jinja2 import Environment, FileSystemLoader
from metadata_extractor import SchemaMetadata, load_schema_metadata

JPK_REQUIRED_FIELDS = {
    "prefixes": [
        # this is a subset of fields that are required.
        # basically, some of the fields that JPK makes mandatory we fill out anyway and so
        # there is no real need to pollute the end file with them. let's stick to the
        # ones that we care about
        "JPK.Deklaracja.PozycjeSzczegolowe",
        "JPK.Ewidencja.SprzedazWiersz",
        "JPK.Ewidencja.ZakupWiersz",
    ],
    "defaults": {
        # consequence of the above is that the parser will tell us about required fields
        # and their types. so we will initialize the JPK document with these zero values
        # which will eventually get overriden. if not, they would stay as zeroes and
        # everyone would be happy.
        "TKwotaC": "0",
        "TKwotaCNieujemna": "0",
        "TKwotowy": "0",
    },
}

TARGET_DIRS = {
    "FA": "../../internal/sei/generators",
    "JPK": "../../internal/invoicesdb/jpk/generators",
}


def filter_required_fields(required_fields: dict[str, str]) -> list[tuple[str, str]]:
    """Filter required fields based on configured prefixes and supported types."""
    return [
        (name, typ)
        for name, typ in required_fields.items()
        if any(name.startswith(prefix) for prefix in JPK_REQUIRED_FIELDS["prefixes"])
        and typ in JPK_REQUIRED_FIELDS["defaults"]
    ]


def generate_go_code(
    variable_prefix: str,
    package_name: str,
    metadata: SchemaMetadata,
) -> str:
    """Generate Go source code using Jinja2 template."""
    env = Environment(loader=FileSystemLoader("./templates"))
    template = env.get_template("go_template.jinja2")

    return template.render(
        package_name=package_name,
        variable_prefix=variable_prefix,
        children_order=metadata.children_order,
        required_defaults=filter_required_fields(metadata.required_fields),
        array_elements=sorted(metadata.array_nodes),
        JPK_REQUIRED_FIELDS=JPK_REQUIRED_FIELDS,
    )


def main():
    schemas_dir = Path(__file__).parent.parent

    for schema_file in sorted(schemas_dir.glob("*.xsd")):
        schema_file_type = schema_file.stem.split("_")[0]
        variable_prefix = schema_file.stem
        package_name = variable_prefix.lower()

        target_dir = Path(TARGET_DIRS[schema_file_type]) / package_name
        target_dir.mkdir(parents=True, exist_ok=True)

        metadata = load_schema_metadata(str(schemas_dir / schema_file.name))

        go_code = generate_go_code(
            variable_prefix=variable_prefix,
            package_name=package_name,
            metadata=metadata,
        )

        target_filename = target_dir / "schema_ordering.go"
        with open(target_filename, "w") as f:
            f.write(go_code)

        print(f"  gofumpt -w {target_filename}")
        call(["gofumpt", "-w", str(target_filename)])


if __name__ == "__main__":
    main()
