using UnrealBuildTool;
using System.Collections.Generic;

public class NeonDriveTarget : TargetRules
{
    public NeonDriveTarget(TargetInfo Target) : base(Target)
    {
        Type = TargetType.Game;
        DefaultBuildSettings = BuildSettingsVersion.V5;
        IncludeOrderVersion = EngineIncludeOrderVersion.Latest;
        ExtraModuleNames.Add("NeonDrive");
    }
}
